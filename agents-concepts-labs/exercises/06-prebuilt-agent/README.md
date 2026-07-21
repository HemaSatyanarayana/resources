# 06 — The prebuilt agent: `ToolNode`, `tools_condition`, `create_agent`

> **Concepts:** `ToolNode` · `tools_condition` · `create_agent` · system prompts · tool error handling · when to stop using the prebuilt

## The idea

You wrote the loop by hand in 05. Now delete it — and notice exactly what you
are deleting.

| Exercise 05, by hand | Exercise 06, prebuilt |
|---|---|
| `run_tools` (~15 lines) | `ToolNode(TOOLS)` |
| `should_continue` (~5 lines) | `tools_condition` |
| the whole `StateGraph` | `create_agent(model, tools)` |

`ToolNode` **is a node** — pass it to `add_node` directly, it needs no wrapper.
`tools_condition` **is your router**, with the same signature and the same two
return values. This is not a coincidence: 05 was a re-implementation of these.

### What `ToolNode` does that your version did not

Your `run_tools` survived an *unknown tool*. `ToolNode` also survives a tool
that **raises**:

```python
tool_call("add", {"a": "two", "b": 1})
# -> ToolMessage("Error invoking tool 'add' with kwargs ... a: Input should be a valid number")
```

Both failures become **observations** — text the model reads on its next turn
and can react to. Nothing unwinds the run. That is the discipline this whole
loop is built on: inside an agent, an error is data, not an exception.

(You can turn this off — `ToolNode(tools, handle_tool_errors=False)` — when you
would rather fail loudly than let a model improvise around a broken dependency.)

### The system prompt goes in the call, not in the state

```python
messages = [SystemMessage(system_prompt), *state["messages"]]   # ✅
state["messages"].insert(0, SystemMessage(system_prompt))       # ❌
```

Put it in state and it gets checkpointed with the thread, re-sent on every turn,
duplicated by careless nodes, and frozen forever — you can never change the
prompt for an existing conversation. Prepend it at call time and the prompt
stays a property of *the code*, while state stays a record of *the
conversation*. `create_agent` does exactly this, which is why a system prompt
never shows up in its output messages.

### So when do you write the graph yourself?

Reach for `create_agent` when you want the standard loop, which is more often
than beginners expect. Drop to a `StateGraph` when you need something the loop
does not have a shape for:

- extra state beyond messages that other nodes read and write (04's `category`)
- steps that are not "call model" or "call tool" — validate, retrieve, classify,
  summarise, escalate
- your own budget and give-up path (03), rather than `recursion_limit` raising
- a human approval gate in the middle (08)
- fan-out to parallel workers (10)

The prebuilt is a preset, not a ceiling. You can also extend it via
`middleware=[...]` without abandoning it — worth reading about once the rest of
this course has landed.

### `ToolNode`, precisely

Your `run_tools` and `ToolNode` do the same job; `ToolNode` does more of it.
Feature by feature, so you know exactly what you are handing over:

| | Your 05 version | `ToolNode` |
|---|---|---|
| Reads `.tool_calls` off the last message | ✅ | ✅ |
| One `ToolMessage` per call, ids preserved | ✅ | ✅ |
| Handles several calls in one message | sequentially | **concurrently** |
| Unknown tool name → error message | ✅ (you wrote it) | ✅ |
| Tool *raises* → error message | ❌ crashes the run | ✅ |
| Injecting graph state into a tool | ❌ | ✅ (`InjectedState`) |

Two of those deserve elaboration.

**Concurrency.** Given three tool calls in one message, `ToolNode` runs them at
the same time rather than one after another — three 200ms API calls take 200ms,
not 600ms. Free latency, with one condition attached: your tools must tolerate
running concurrently. Pure functions and read-only lookups are fine. Two tools
mutating the same row are not, and the framework cannot know the difference.

**Tool exceptions.** This is the case your version did not cover:

```python
tool_call("add", {"a": "two", "b": 1})
# -> ToolMessage("Error invoking tool 'add' with kwargs {...}: a: Input should be a valid number")
```

The `@tool` decorator validates arguments against the schema it generated, so a
model that sends a string where a number belongs fails *before* your function
body runs. `ToolNode` catches it and converts it into an observation, exactly as
exercise 05 argued it should.

That behaviour is a flag, not a law:

```python
ToolNode(tools, handle_tool_errors=False)   # let it propagate
```

Turn it off when a failing tool means a broken dependency rather than a
correctable mistake. With it on, a database outage looks to the model like a bad
answer — and a helpful model will route around it by inventing something. Silent
plausible wrongness is worse than a crash. The rule of thumb: **on** for errors
the model caused and can fix, **off** for errors it cannot.

### `tools_condition` is your router, verbatim

```python
def tools_condition(state) -> Literal["tools", "__end__"]:
    # last message has tool_calls -> "tools", else END
```

Same signature, same two return values, same logic as the `should_continue` you
wrote. Pass it to `add_conditional_edges`; never call it yourself.

It expects the tool node to be **named `"tools"`** — that string is baked in. If
you name your node `"my_tools"`, pass a `path_map` (exercise 03) or the routing
silently points nowhere. This is the small tax of prebuilts: they encode
conventions you must now match.

### What `create_agent` actually builds

```python
create_agent(model, TOOLS, system_prompt=SYSTEM_PROMPT)
```

It returns a compiled graph with nodes named `model` and `tools`, wired
`START → model → (tools_condition) → tools | END`, `tools → model` — the graph
you just built by hand, which is why the parity test can compare them
message-for-message.

The important thing is what it is *not*: not a special "Agent" class, not a
different runtime, not something that has to be used all-or-nothing. It is a
`CompiledStateGraph` like any other, so `.invoke`, `.stream`, `checkpointer=`,
`.get_state()` and everything in exercises 07–09 apply unchanged.

It also takes arguments that preview the rest of the course — worth recognising
now, even before you use them:

| Argument | What it does | Covered in |
|---|---|---|
| `checkpointer=` | conversation memory | 07 |
| `interrupt_before=` | pause before a node | 08 |
| `response_format=` | structured output instead of prose | 04's note |
| `middleware=` | hooks around the loop without abandoning the prebuilt | — |
| `state_schema=` | extra state keys beyond `messages` | 04 |

### System prompt: in the call, never in the state

```python
messages = [SystemMessage(system_prompt), *state["messages"]]   # ✅
state["messages"].insert(0, SystemMessage(system_prompt))       # ❌
```

Exercise 04 made this point; here is the full cost of getting it wrong, because
it is one of the most common design mistakes in agent code.

Put the system prompt in `messages` and it is **checkpointed** (07), so it
becomes part of the conversation's permanent record. It is **re-sent** on every
turn regardless — no saving. It is **duplicated** by any node that prepends
again, and now the model receives two conflicting briefs. Worst of all it is
**frozen**: every conversation started before your prompt fix keeps the old
wording forever, and you cannot correct it without rewriting stored history.

Prepend at call time and the split is clean: the prompt is a property of *the
code*, state is a record of *the conversation*. Deploy a better prompt and every
thread — including ones mid-flight — picks it up on the next turn.

This is why `create_agent` takes `system_prompt` as a parameter and why its
output messages never contain a `SystemMessage`. Take yours as a parameter with
a default too; one test builds a pirate, and more importantly, a prompt you can
vary is a prompt you can A/B test.

### When to stop using the prebuilt

Reach for `create_agent` when you want the standard loop — more often than
beginners expect, and the default answer for "an agent with tools".

Drop to a `StateGraph` when your problem is not shaped like that loop:

- **Extra state** that other nodes read and write (04's `category`)
- **Steps that are not "call model" or "call tool"** — validate, retrieve,
  classify, summarise, escalate, format
- **Your own budget and give-up path** (03) instead of `recursion_limit` raising
- **A human approval gate** in the middle (08)
- **Fan-out to parallel workers** (10)
- **Multiple models** — a cheap one to route, an expensive one to answer

The honest framing: the prebuilt is a **preset**, not a ceiling, and the two are
not a one-way door. You can start with `create_agent`, and the day it stops
fitting, you already know what it was doing — you wrote it in exercise 05.
That is the entire reason this exercise comes second.

## The imports you need

```python
from __future__ import annotations

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import SystemMessage
from langgraph.graph import START, MessagesState, StateGraph
from langgraph.prebuilt import ToolNode, tools_condition

from labgraph import add, multiply, search
```

New here, and note the new module path — `langgraph.prebuilt`, not
`langgraph.graph`. That separation is meaningful: `langgraph.graph` is the
primitives (state, nodes, edges), `langgraph.prebuilt` is opinionated
assemblies built *out of* those primitives.

**`ToolNode`** — a ready-made replacement for your `run_tools`. It **is a node**:
pass it straight to `add_node`, no wrapper, no lambda. Constructed with the tool
list, `ToolNode(TOOLS)`, it reads `.tool_calls` off the last message, runs them
(concurrently when there are several), and returns the `ToolMessage`s — plus it
catches exceptions and hands the model an error string instead of crashing.

**`tools_condition`** — a ready-made replacement for your `should_continue`,
with the same signature and the same two return values. Pass it directly to
`add_conditional_edges`; you never call it yourself.

Notice `END` is no longer imported. `tools_condition` returns `END` internally,
and your `build_graph` never names it — the only edges you write explicitly are
`START -> model` and `tools -> model`.

Inside `build_prebuilt_agent` you need one more, and it comes from **LangChain**,
not LangGraph:

```python
from langchain.agents import create_agent
```

Import it inside the function rather than at module top if you like — it pulls
in a larger dependency tree, and keeping it local makes the point that the rest
of the file does not depend on it.

## What to build

### 1. `TOOLS` and `SYSTEM_PROMPT`

```python
TOOLS = [add, multiply, search]

SYSTEM_PROMPT = (
    "You are a careful assistant. Use your tools for arithmetic and for "
    "questions about LangGraph. Never guess at a calculation."
)
```

The wording is yours; the spec compares against your own constant.

### 2. `build_graph(model, system_prompt=SYSTEM_PROMPT)`

Exercise 05's graph, rebuilt on the prebuilt parts. Two swaps:

- `run_tools` → `ToolNode(TOOLS)`
- `should_continue` → `tools_condition`

The model node is still yours, because that is where the system prompt goes:

```python
messages = [SystemMessage(system_prompt), *state["messages"]]
```

Prepend it at call time. Do **not** append it to state — a system message in
`messages` gets checkpointed, re-sent, and duplicated on every turn, and you can
never change it afterwards.

Wiring is unchanged from 05:

```
START -> model
model -> (tools_condition) -> tools | END
tools -> model
```

Name the nodes `model` and `tools` — what `create_agent` calls them, so the
parity test compares like with like.

Take `system_prompt` as a **parameter with a default**, not a hardcoded
constant: one test builds a pirate.

### 3. `build_prebuilt_agent(model, system_prompt=SYSTEM_PROMPT)`

```python
return create_agent(model, TOOLS, system_prompt=system_prompt)
```

It returns a compiled graph — `.invoke`, `.stream`, checkpointers, all of it
works exactly as on a graph you built. It is not a different kind of object.

Yes, that is the whole implementation. The exercise is not typing it; the
exercise is `test_handwritten_and_prebuilt_agree_message_for_message` passing,
which is only interesting once you have written the other one.

### 4. Optional: a runnable demo

```python
if __name__ == "__main__":
    from labgraph import get_chat_model, print_messages

    question = {"messages": [("user", "What is 128 times 4?")]}

    for label, factory in (("hand-wired", build_graph), ("prebuilt", build_prebuilt_agent)):
        print(f"\n== {label} ==")
        print_messages(factory(get_chat_model()).invoke(question)["messages"])
```

## Run it

```bash
uv run pytest exercises/06-prebuilt-agent -v
uv run python exercises/06-prebuilt-agent/graph.py     # needs .env
```

## Think about it

- `ToolNode` runs parallel tool calls concurrently. What breaks if two of your
  tools write to the same row of a database?
- `handle_tool_errors=True` means a broken database looks to the model like a
  bad answer. When would you rather the run just crashed?
- The system prompt is not in state — so what happens to it when you resume a
  conversation from a checkpoint tomorrow with different code? (Exercise 07.)
- `create_agent` accepts `checkpointer=` and `response_format=`. Guess what each
  one does before you look them up.
- Your `TOOLS` list is fixed at build time. How would you give one user access
  to a tool another user must not have?

<details>
<summary>Reference solution</summary>

```python
from __future__ import annotations

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import SystemMessage
from langgraph.graph import START, MessagesState, StateGraph
from langgraph.prebuilt import ToolNode, tools_condition

from labgraph import add, multiply, search

TOOLS = [add, multiply, search]

SYSTEM_PROMPT = (
    "You are a careful assistant. Use your tools for arithmetic and for "
    "questions about LangGraph. Never guess at a calculation."
)


def build_graph(model: BaseChatModel, system_prompt: str = SYSTEM_PROMPT):
    bound = model.bind_tools(TOOLS)

    def call_model(state: MessagesState) -> dict:
        messages = [SystemMessage(system_prompt), *state["messages"]]
        return {"messages": [bound.invoke(messages)]}

    builder = StateGraph(MessagesState)

    builder.add_node("model", call_model)
    builder.add_node("tools", ToolNode(TOOLS))

    builder.add_edge(START, "model")
    builder.add_conditional_edges("model", tools_condition)
    builder.add_edge("tools", "model")

    return builder.compile()


def build_prebuilt_agent(model: BaseChatModel, system_prompt: str = SYSTEM_PROMPT):
    from langchain.agents import create_agent

    return create_agent(model, TOOLS, system_prompt=system_prompt)
```

Note there is no `add_conditional_edges("tools", ...)` — `tools -> model` is an
unconditional edge. Only one node in a cycle needs to make a decision.

</details>

---

Next: [07 — Memory](../07-memory) — the agent currently forgets everything the
moment `.invoke()` returns.
