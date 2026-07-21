# 05 — Tool calling from scratch: the ReAct loop

> **Concepts:** `bind_tools` · `AIMessage.tool_calls` · `ToolMessage` · `tool_call_id` · parallel tool calls · the agent loop

## The idea

Here is the thing worth internalising before you write a line of code:

> **The model never runs a tool.** It emits a *request* — a name and some JSON
> arguments — and stops. Your code runs the tool and hands the result back.

That is why an agent is a **loop** and not a function call. Four messages, one
cycle:

| # | Message | Who made it |
|---|---|---|
| 1 | `HumanMessage("what is 2+3?")` | the user |
| 2 | `AIMessage(content="", tool_calls=[{"name": "add", "args": {...}, "id": "call_abc"}])` | the model |
| 3 | `ToolMessage(content="5.0", tool_call_id="call_abc")` | **you** |
| 4 | `AIMessage("2 + 3 is 5.")` | the model, having seen #3 |

Steps 2→3→4 are the loop. Take away the `tools -> model` edge and you compute
step 3 and throw it away — the model answers step 4 having never seen the
result. That single missing edge is the most common broken-agent bug.

### `bind_tools` is what makes step 2 possible

```python
bound = model.bind_tools([add, multiply, search])
```

This returns a *new* model that ships the tool schemas — names, descriptions,
argument types, pulled from your `@tool` function's signature and docstring —
with every request. Unbound, the model does not know the tools exist and will
cheerfully do arithmetic in its head and get it wrong.

**Your docstrings are prompt.** `search`'s docstring says *"Use for questions
about how LangGraph works"* — that sentence is the only thing telling the model
when to reach for it. A vague docstring is a tool the model misuses.

### `tool_call_id` is load-bearing

Every result must carry the id of the request it answers. Providers validate
this: a `ToolMessage` with a missing or unmatched id is a hard API error on the
next call, not a silent degradation. And since a model may request **several
tools at once** — `.tool_calls` is a list — the id is the only thing keeping the
answers straight.

### Never raise on a bad tool call

A model can ask for a tool that does not exist, or pass arguments that blow up
the tool. Two options:

```python
raise KeyError(name)                       # run over, stack trace to the user
ToolMessage(f"Error: unknown tool {name}") # model sees it and tries again
```

The second is almost always right. The error becomes an *observation*, and a
decent model recovers from it. This is a distinctive property of agents: the
error-handling path feeds back into the reasoning instead of unwinding it.

### What stops the loop?

Nothing in this graph, deliberately. The model stops asking for tools, or
LangGraph's `recursion_limit` (default 25) raises `GraphRecursionError`. Compare
that to exercise 03, where *you* owned the budget in state.

For a toy this is fine. In production it usually is not — a limit that raises
gives your user a stack trace instead of an answer. The mature version keeps a
counter in state and, when it runs out, exits through a node that says
"I couldn't finish this" in the model's voice. You now know how to build that;
it is exercise 03's pattern applied here.

### The full round trip, one HTTP request at a time

The message table above shows the transcript. Here is what actually crosses the
wire, because "the model calls a tool" hides the only part that matters.

**Before anything: `bind_tools` builds schemas.** When you call
`model.bind_tools([add, multiply, search])`, LangChain converts each `@tool`
function into a JSON Schema by reading its signature and docstring:

```json
{"name": "add",
 "description": "Add two numbers together.",
 "parameters": {"type": "object",
                "properties": {"a": {"type": "number"},
                               "b": {"type": "number"}},
                "required": ["a", "b"]}}
```

Note what became the `description`: **your docstring**. And what became the
parameter types: **your type hints**. Both are shipped to the provider on every
request, and they are the model's only information about when and how to use the
tool. A tool whose docstring says "does the thing" is a tool the model will
misuse, and no amount of system-prompt tuning fixes it.

**Request 1** carries the conversation *plus* those schemas. The model reads
"what is 2+3?", sees a tool called `add` described as adding numbers, and
decides to use it.

**Response 1** is the part that surprises people. The model does not return an
answer; it returns a **request**, and stops:

```python
AIMessage(content="", tool_calls=[
    {"name": "add", "args": {"a": 2, "b": 3}, "id": "call_abc123"}
])
```

`content` is empty. The payload is in `.tool_calls` — a name it generated, an
arguments dict it generated, and an id **it** minted to label this specific
request.

**Your code runs the tool.** Nothing else can: the model is a text API on
someone else's servers. It has no filesystem, no network to your database, no
access to `LEDGER`. This is not a limitation to work around, it is the entire
security boundary of agents — *you* decide what a tool call is allowed to do,
because *you* are the one executing it. Exercise 08 is that sentence taken
seriously.

**Request 2** is the whole conversation again, now with your `ToolMessage`
appended. The model finally sees `"5.0"` and writes prose.

Two calls, four messages, one tool. Now re-read the graph: `model → tools →
model` is exactly that sequence, and the `tools → model` edge is request 2.
Delete it and you have computed the answer and thrown it away.

### Why `tool_call_id` is not optional

Every `ToolMessage` must carry the id of the `tool_call` it answers. Providers
**validate** this server-side, and the failure mode is unhelpful: an API error on
the *next* request, complaining about message structure, with no hint about which
result was mismatched.

The id exists because tool calls are not necessarily one at a time. A model can
emit several in a single `AIMessage`:

```python
AIMessage(content="", tool_calls=[
    {"name": "add",      "args": {...}, "id": "call_1"},
    {"name": "multiply", "args": {...}, "id": "call_2"},
])
```

This is **parallel tool calling**, and it is why `run_tools` returns a *list* of
ToolMessages rather than one. Position is not the pairing mechanism — the id is.
Return one result for a two-call message and the conversation is malformed;
return them with swapped ids and the model reads the answers backwards, which is
worse, because nothing errors.

### Errors as observations: the defining habit

Two ways to handle a tool call for a tool that does not exist:

```python
raise KeyError(name)                        # the run is over
ToolMessage(f"Error: unknown tool {name}")  # the model reads it and adapts
```

Ordinary programs prefer the first: fail fast, surface the stack trace, let the
caller decide. Agents usually want the second, and the reason is structural — an
agent has a **feedback loop** that ordinary code does not. The error goes back
into the conversation, the model sees what it did wrong, and it gets another
turn to try something else. A hallucinated tool name becomes a recoverable
detour instead of a 500.

This inverts a habit worth being conscious about. You are deliberately *not*
letting the exception propagate, because propagating it destroys the one thing
that could fix it. The judgement call is which failures deserve this treatment:
a bad argument or a wrong tool name, yes — the model can correct those. A
database that is down, no; the model will "recover" by inventing an answer, and
you have converted an outage into a plausible lie. Exercise 06's
`handle_tool_errors` is where that choice becomes a flag.

### What stops the loop — and why this version is not production-ready

Look at the graph honestly: there is **no budget**. It terminates when the model
stops asking for tools. Compare exercise 03, where you owned an explicit
`MAX_REVISIONS` counter in state.

The safety net is `recursion_limit` (exercise 03) — and a net is all it is.
`GraphRecursionError` reaches your user as a stack trace, after you have already
paid for every call it took to get there.

The mature version keeps a counter in state, and when it runs out, routes to a
node that produces a *graceful* answer in the model's voice: "I wasn't able to
finish this — here is what I found." You already know how to build that; it is
exercise 03's budget pattern applied to exercise 05's loop, and combining them
unprompted is a good sign you have both.

### The cost curve nobody mentions

Each cycle appends at least two messages, and **the entire conversation is
re-sent every time**. Turn *n* costs roughly proportional to *n*, so a 20-step
agent run is not 20 times the cost of one step — it is closer to the sum
1+2+…+20, an order of magnitude worse.

Three consequences, all of which become real work later:

- **Latency compounds.** Every turn re-processes everything before it.
- **The context window is finite.** Long tool outputs (a 50KB API response, a
  full HTML page) fill it fastest, which is why real tools return *summaries*,
  not dumps.
- **Signal dilutes.** More irrelevant history means worse answers, not just
  pricier ones.

The standard mitigations are trimming, summarising, and moving detail out to a
store the agent can query. All three are ordinary nodes operating on
`state["messages"]` — you have had the tools to write them since exercise 02.

## The imports you need

```python
from __future__ import annotations

from typing import Literal

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import ToolMessage
from langgraph.graph import END, START, MessagesState, StateGraph

from labgraph import add, multiply, search
```

Carried over from 04: `Literal`, `BaseChatModel`, `MessagesState`,
`StateGraph`, `START`, `END`. New here:

**`ToolMessage`** — the fourth message type, and the one you now produce
yourself. The transcript so far has been `SystemMessage` (instructions),
`HumanMessage` (user), `AIMessage` (model). A `ToolMessage` is **the result of a
tool**, and it exists as its own type because the model must be able to tell
"here is what your tool returned" apart from "here is what a person said". It
takes three things that matter:

```python
ToolMessage(content="5.0", tool_call_id="call_abc123", name="add")
```

`content` must be **text** — tool results travel as strings. `tool_call_id` is
mandatory and links the result to the request that asked for it. `name` is the
tool's name, and is conventional rather than load-bearing.

**`add`, `multiply`, `search` from `labgraph`** — the toy tools, already
written. Read `labgraph/tools.py`; each is ten lines. What matters is that they
are decorated with LangChain's `@tool`, which turns a plain function into a
`BaseTool`: it reads the signature and docstring and generates a **JSON schema**
that gets sent to the model. Two consequences:

- `tool.name` is the function name — the string the model will send back.
- `tool.invoke({"a": 2, "b": 3})` calls it with a dict of arguments, not
  positionally.
- **Your docstring is prompt.** `search`'s docstring says *"Use for questions
  about how LangGraph works"* — that sentence is the entire basis on which the
  model decides to reach for it.

You do **not** import `@tool` in this exercise (the tools already exist), but you
will write one in exercise 08.

Two things you will use that need no import: `model.bind_tools(TOOLS)`, and
`AIMessage.tool_calls` — a plain list of dicts hanging off the reply.

## What to build

### 1. `TOOLS` and a name lookup

```python
TOOLS = [add, multiply, search]
TOOLS_BY_NAME = {t.name: t for t in TOOLS}
```

A model asks for a tool **by name** — a string it generated — so you need a
lookup, and you need to survive a name that is not in it.

### 2. `make_model_node(model)` — a node factory

Returns a node that calls the model **with the tools attached**.

`model.bind_tools(TOOLS)` returns a *new* model that advertises those tools in
every request. Without it the model has no idea the tools exist and will happily
invent an answer to "what is 128 × 4".

Bind **once**, out in the factory — not inside the node. The node runs on every
turn of the loop; binding is setup. (The spec asserts this.)

The node itself is the same shape as exercise 04: call the bound model with
`state["messages"]`, return `{"messages": [reply]}`.

### 3. `run_tools(state) -> dict`

Execute every tool call on the last message and return the results.

The last message is an `AIMessage` whose `.tool_calls` is a list of dicts:

```python
{"name": "add", "args": {"a": 2, "b": 3}, "id": "call_abc123"}
```

For each one:

1. Look the tool up in `TOOLS_BY_NAME`.
2. If it is missing, do **not** raise. Build a `ToolMessage` whose content says
   `unknown tool` and names it. A raise ends the run; an error message goes back
   to the model, which can then try something else. Models hallucinate tool
   names — this path is load-bearing.
3. Otherwise call `tool.invoke(call["args"])` and **stringify** the result.
   `add` returns a float; the wire wants text.
4. Wrap it in `ToolMessage(content=..., tool_call_id=call["id"], name=...)`.

A model can ask for several tools at once, so return a **list** of ToolMessages
— one per call, in the same order.

### 4. `should_continue(state) -> Literal["tools", "__end__"]`

The router: did the model ask for a tool? Look at the last message; if it has a
non-empty `tool_calls`, return `"tools"`, otherwise `END`.

Use `getattr(last, "tool_calls", None)` — the last message is not always an
`AIMessage`, and only `AIMessage` has that attribute.

Note what this router does **not** do: call a model, or decide *which* tool. The
model already decided. The router only reads state.

### 5. `build_graph(model)`

```
START  -> model
model  -> (conditional: should_continue) -> tools | END
tools  -> model          <- the edge that makes it an agent
```

Node names: `model`, `tools`.

That `tools -> model` edge is the one people forget. Without it the tool result
is computed and then thrown away, never seen by the model.

Notice there is no step budget here — unlike exercise 03. A well-behaved model
stops on its own, and `recursion_limit` catches one that does not.

### 6. Optional: a runnable demo

```python
if __name__ == "__main__":
    from labgraph import get_chat_model, print_messages

    graph = build_graph(get_chat_model())
    result = graph.invoke(
        {"messages": [("user", "What is 128 times 4? Then tell me what a reducer is.")]}
    )

    print_messages(result["messages"])
```

## Run it

```bash
uv run pytest exercises/05-tool-calling -v
uv run python exercises/05-tool-calling/graph.py     # needs .env
```

## Think about it

- The test `test_the_model_sees_the_tool_result_on_the_second_turn` asserts on
  the model's *prompt*, not the graph output. Why is that the test that would
  actually catch a missing `tools -> model` edge?
- The conversation grows by 2+ messages per cycle. What happens to cost and
  latency on turn 20? What would you drop, and what must you never drop?
- `search` is read-only. What would you change about this loop if a tool
  charged a credit card? (Exercise 08.)
- If a tool takes 30 seconds, the whole graph blocks. Where would you put the
  concurrency? (Exercise 10.)
- Nothing validates the model's `args` against the tool's signature. What
  happens when it sends `{"a": "two"}`?

<details>
<summary>Reference solution</summary>

```python
from __future__ import annotations

from typing import Literal

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import ToolMessage
from langgraph.graph import END, START, MessagesState, StateGraph

from labgraph import add, multiply, search

TOOLS = [add, multiply, search]
TOOLS_BY_NAME = {t.name: t for t in TOOLS}


def make_model_node(model: BaseChatModel):
    bound = model.bind_tools(TOOLS)

    def call_model(state: MessagesState) -> dict:
        return {"messages": [bound.invoke(state["messages"])]}

    return call_model


def run_tools(state: MessagesState) -> dict:
    last = state["messages"][-1]
    results = []

    for call in last.tool_calls:
        tool = TOOLS_BY_NAME.get(call["name"])
        if tool is None:
            content = f"Error: unknown tool {call['name']!r}"
        else:
            content = str(tool.invoke(call["args"]))
        results.append(
            ToolMessage(content=content, tool_call_id=call["id"], name=call["name"])
        )

    return {"messages": results}


def should_continue(state: MessagesState) -> Literal["tools", "__end__"]:
    last = state["messages"][-1]
    if getattr(last, "tool_calls", None):
        return "tools"
    return END


def build_graph(model: BaseChatModel):
    builder = StateGraph(MessagesState)

    builder.add_node("model", make_model_node(model))
    builder.add_node("tools", run_tools)

    builder.add_edge(START, "model")
    builder.add_conditional_edges("model", should_continue)
    builder.add_edge("tools", "model")

    return builder.compile()
```

`tool.invoke(call)` — passing the whole call dict rather than just `call["args"]`
— returns a fully-formed `ToolMessage` with the id already set. The long form
above is written out so you can see what that shortcut does for you.

</details>

---

Next: [06 — The prebuilt agent](../06-prebuilt-agent) — throw most of this away,
and know exactly what you threw.
