# 04 — An LLM inside a node

> **Concepts:** `MessagesState` · model injection · node factories · parsing LLM output · defensive fallbacks

## The idea

Nothing structural changes when a node calls a model. It is still
`state -> partial update`. What changes is that the update is now
**non-deterministic**, and that has two consequences you should build habits
around right now.

### 1. Inject the model, never construct it inside

```python
def build_graph(model: BaseChatModel):   # yes
def build_graph():                        # no — where does the model come from?
    model = ChatOpenAI(...)
```

Because the model is an argument, the exact same graph runs against your real
provider in `__main__` and against `ScriptedChatModel` in the tests. Every test
in this exercise is deterministic and needs no API key — that is not a testing
trick, it is the direct payoff of dependency injection. A graph that builds its
own model can only be tested by calling a real one.

### 2. LLMs return prose; graphs route on values

The classifier asks for one word. It will still occasionally return
`"I think this is a billing issue!"`. So the node **normalises and validates**:

```python
label = reply.content.strip().lower().strip(".!?\"' ")
if label not in CATEGORIES:
    label = "general"          # never route somewhere that does not exist
```

Skipping that fallback is how you get a `KeyError` in production at 3am. Treat
every model output as untrusted input — because it is.

> Later you will meet `with_structured_output()` and `response_format`, which
> get the model to emit a schema-validated object instead of prose. They make
> parsing tidier; they do not remove the need for a fallback.

### `MessagesState` and node factories

`MessagesState` is the prebuilt version of what you wrote in exercise 02 —
`{"messages": Annotated[list, add_messages]}`. Subclass it to add keys:

```python
class SupportState(MessagesState):
    category: str
```

A **node factory** is a closure that returns a node function:

```python
def make_responder(model, category):
    def respond(state):
        ...
    return respond
```

This is how you parameterise nodes without globals, and it lets you register
three personas in a loop instead of writing three near-identical functions.

### What belongs in `messages`?

The classifier's one-word reply is *scaffolding*, not conversation, so it
returns `{"category": ...}` and leaves `messages` alone. The responder's reply
*is* conversation, so it returns `{"messages": [reply]}`. Being deliberate about
this matters: everything you put in `messages` gets re-sent to the model on
every later turn, costing tokens and inviting confusion.

Also note the system prompt is **prepended at call time**, not stored in state:

```python
model.invoke([SystemMessage(PERSONAS[category]), *state["messages"]])
```

Keeping it out of state means each node can use its own persona over the same
shared history.

```
                        ┌──▶ billing ───┐
START ──▶ classify ──┬──┼──▶ technical ─┼──▶ END
                     └──┴──▶ general ───┘
```

### What a chat model call actually is

Before the graph parts, be clear about the thing inside the node, because a lot
of confusion later comes from a fuzzy model here.

A chat model is a **stateless function from a list of messages to one message.**
That is the entire interface:

```python
reply = model.invoke([SystemMessage("..."), HumanMessage("...")])
# reply is an AIMessage
```

It remembers nothing. Every notion of "the conversation" is something *you*
reconstruct by sending the whole history every single time. When exercise 07
adds memory, it is not teaching the model to remember — it is automating the
bookkeeping of what to re-send.

The four message types are not interchangeable wrappers around strings. Each
maps to a **role** the provider treats differently:

| Type | Role | Who writes it | What it is for |
|---|---|---|---|
| `SystemMessage` | system | you | standing instructions: persona, rules, format |
| `HumanMessage` | user | the user | the request |
| `AIMessage` | assistant | the model | its reply — and, from 05, its tool requests |
| `ToolMessage` | tool | **your code** | the result of running a tool (exercise 05) |

Models are trained to weight system content as instruction and user content as
request. Putting "you are a billing specialist" in a `HumanMessage` measurably
degrades adherence — it reads as something the user claimed, not a rule.

`reply.content` is typed `str | list`: providers that return multi-part content
(text plus images, or reasoning blocks) use the list form. That is why the
reference solution writes `str(reply.content)` before parsing rather than
assuming a string.

### Non-determinism changes what you can rely on

The structural shape of the node has not changed — `state -> update`, same as
exercise 01. What changed is that **the same input can produce a different
output**, and several habits you have from ordinary code stop working:

- **You cannot assert on exact output.** Not in tests, not in downstream logic.
  You assert on *structure* ("it is one of three categories") and let the text
  vary.
- **You cannot assume the format you asked for.** "Reply with the single word
  only" is a request, not a constraint. Sometimes you get `"Billing."`,
  sometimes `"I think this is a billing issue!"`.
- **Retrying is not free and not idempotent.** A retried node calls the model
  again: new cost, new latency, possibly a different answer.
- **Every call costs money and time.** This graph makes *two* calls per request
  — one to classify, one to answer. That is a real design decision (see *Think
  about it*), not an implementation detail.

### Dependency injection, in more detail

```python
def build_graph(model: BaseChatModel):        # yes
def build_graph():                            # no
    model = ChatOpenAI(...)
```

The second version welds three unrelated decisions together: *what the graph
does*, *which provider you use*, and *where credentials come from*. Every one of
those changes independently in real life.

Passing the model as an argument is what makes all of this possible at once:

| | With injection | Without |
|---|---|---|
| **Tests** | `ScriptedChatModel` — instant, free, deterministic | a real API call per test |
| **Provider switch** | change one line in `__main__` | edit the graph |
| **Two models in one graph** | pass both; cheap one to classify, strong one to answer | impossible without surgery |
| **Rate limits / retries / caching** | wrap the model, graph untouched | tangled into node code |

The tests in this exercise are the proof: every one runs offline, in
milliseconds, with no key — and they exercise *your real graph*, not a mock of
it. That is not a testing trick. It is the direct, mechanical payoff of not
constructing your dependencies inside the thing that uses them.

### Parsing model output: the ladder

The classifier's job is to turn prose into a value the router can branch on.
Work down the ladder, and never skip the last rung:

```python
raw = str(reply.content)          # 1. it may not be a str
label = raw.strip()               # 2. whitespace
label = label.lower()             # 3. case: "Billing" == "billing"
label = label.strip(".!?\"' ")    # 4. punctuation: "Billing." -> "billing"
if label not in CATEGORIES:       # 5. THE IMPORTANT ONE
    label = "general"
```

Steps 2–4 handle a model that basically complied. Step 5 handles one that did
not, and it is the difference between a graph that degrades and one that
crashes. Without it, a chatty reply becomes a `category` naming no node, and
LangGraph fails at run time — in production, at 3am, on the one input you never
tried.

Choosing the *fallback* is a real decision. `"general"` is right here because a
general agent can handle a billing question adequately, while a billing
specialist handling a technical question is worse. When there is no safe
default, the honest move is an explicit "I could not classify this" branch that
asks the user — not a coin flip.

> **Where this goes later.** `with_structured_output()` and `create_agent`'s
> `response_format` make the model emit a schema-validated object rather than
> prose, which removes steps 1–4. They do **not** remove step 5: the model can
> still fail to produce anything valid, and the call can still fail. Structured
> output changes the shape of the problem, not its existence.

### Node factories: closures as configuration

```python
def make_responder(model, category):
    def respond(state):
        reply = model.invoke([SystemMessage(PERSONAS[category]), *state["messages"]])
        return {"messages": [reply]}
    return respond
```

`make_responder` is not a node. It is a **factory** that *returns* one. The
returned `respond` closes over `model` and `category`, so it has everything it
needs while still matching the `(state) -> dict` signature LangGraph requires.

This is the standard answer to "my node needs configuration". The alternatives
are worse: module-level globals (untestable, unsafe with concurrent runs), or
stuffing config into state (it is not state — it never changes during a run, and
it would get checkpointed for no reason).

It also lets you register the three personas in a loop:

```python
for category in CATEGORIES:
    builder.add_node(category, make_responder(model, category))
```

Three near-identical node functions would be the alternative. One classic
gotcha, worth knowing because it does not apply here and people often think it
does: closing over a loop variable *directly* is a late-binding bug in Python —
all closures would see the final value. Passing `category` as a **function
argument** to the factory, as above, binds it immediately and sidesteps the
problem entirely.

### Prompt assembly: what is sent, versus what is stored

Each responder builds its prompt fresh:

```python
[SystemMessage(PERSONAS[category]), *state["messages"]]
```

The system message is **prepended at call time and never stored**. Two payoffs:

1. **Three personas, one history.** Each node applies its own instructions over
   the same shared conversation. If the persona lived in `messages`, the
   transcript would accumulate contradictory system messages as the run
   progressed.
2. **The prompt stays a property of the code.** Change `PERSONAS` and every
   conversation — including ones already in progress — uses the new wording
   immediately. That is usually what you want; exercise 07's *Think about it*
   explores when it is not.

The mirror-image question is what you *do* store, and the classifier answers it:
its one-word reply is **scaffolding, not conversation**, so the node returns
`{"category": label}` and leaves `messages` untouched.

Be deliberate about this every time, because `messages` is the one part of state
you re-send on every subsequent call. Everything in it costs tokens forever, and
irrelevant content does not merely cost money — it degrades answers, by giving
the model more things to pay attention to than the task requires.

## The imports you need

```python
from __future__ import annotations

from typing import Literal, TypedDict

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import SystemMessage
from langgraph.graph import END, START, MessagesState, StateGraph
```

Carried over: `Literal` (03), `END`/`START`/`StateGraph` (01). New here:

**`BaseChatModel`** — the abstract base class every LangChain chat model
inherits from, real or fake. You need it only as a **type annotation**:

```python
def build_graph(model: BaseChatModel):
```

That annotation is the injection contract made explicit. It says "any chat
model", which is precisely why `ChatOpenAI` in `__main__` and
`ScriptedChatModel` in the tests are both acceptable. You never *construct* a
`BaseChatModel` — you receive one.

**`SystemMessage`** — the message type carrying instructions rather than
conversation. Where `HumanMessage` and `AIMessage` (02) are the transcript, a
`SystemMessage` is the standing brief: "you are a billing specialist". Models
weight it differently from user text, which is why it is its own class rather
than a `HumanMessage` with special wording.

**`MessagesState`** — LangGraph's prebuilt state schema, exactly the
`{"messages": Annotated[list, add_messages]}` you hand-wrote in 02. Import it
from `langgraph.graph`. **Subclass it** to add your own keys:

```python
class SupportState(MessagesState):
    category: str
```

It is a `TypedDict` subclass, so subclassing means "all of its keys, plus mine",
and `messages` keeps the `add_messages` reducer. `TypedDict` is imported here
only if you prefer to spell the schema out yourself.

**`.invoke(messages)`** is the one model API you need. Hand it a list of
messages, get one `AIMessage` back. Its `.content` is the text — and it is
`str | list`, which is why the reference solution wraps it in `str(...)` before
parsing.

## What to build

### 1. Three module-level constants

```python
CATEGORIES = ("billing", "technical", "general")

CLASSIFIER_PROMPT = (
    "Classify the user's support request as exactly one of: "
    "billing, technical, general. Reply with the single word only."
)

PERSONAS = {
    "billing": "You are a billing specialist. Be precise about money.",
    "technical": "You are a support engineer. Be concrete and practical.",
    "general": "You are a friendly support agent.",
}
```

The exact prompt wording is yours to change — the spec compares against your own
constant. `PERSONAS` must have one distinct entry per category.

### 2. `SupportState`

`MessagesState` plus a `category: str` key.

### 3. `make_classifier(model)` — a node factory

Returns a **node function** that classifies the latest user message. The node
should:

1. Call the model with `SystemMessage(CLASSIFIER_PROMPT)` followed by the
   conversation so far (`state["messages"]`).
2. Normalise the reply into one of `CATEGORIES`: lowercase it, strip whitespace
   and punctuation (`"Billing."` → `"billing"`). If the result is not a known
   category, fall back to `"general"`.
3. Return `{"category": <category>}`.

Note what it does **not** return: the classifier's reply is scaffolding, not
conversation, so keep it out of `messages`.

### 4. `make_responder(model, category)` — a node factory

Returns a node that calls the model with `SystemMessage(PERSONAS[category])`
followed by `state["messages"]`, and returns the reply as
`{"messages": [reply]}` — a **list**, because `add_messages` appends it.

### 5. `route_by_category(state)`

Return the node named by `state["category"]`. Annotate the return as
`Literal["billing", "technical", "general"]` so the graph draws correctly.

### 6. `build_graph(model)`

```
START -> classify
classify -> (conditional) -> billing | technical | general
each responder -> END
```

Name the classifier node `classify` and each responder after its category. Build
the responders in a **loop over `CATEGORIES`** — three near-identical `add_node`
calls is a smell the factory already solved.

### 7. Optional: a runnable demo

```python
if __name__ == "__main__":
    from labgraph import get_chat_model, print_messages

    graph = build_graph(get_chat_model())
    result = graph.invoke({"messages": [("user", "I was charged twice this month")]})

    print(f"category: {result['category']}")
    print_messages(result["messages"])
```

## Run it

```bash
uv run pytest exercises/04-llm-node -v          # offline, fake model
cp .env.example .env                            # then fill it in
uv run python exercises/04-llm-node/graph.py    # against a real provider
```

The live test at the bottom of the spec skips itself until `.env` is configured.

## Think about it

- The classifier costs a whole model call to produce one word. When is a cheap
  fast model for classification and an expensive one for the answer worth it —
  and how would you wire two different models into this graph?
- What if the user's request is *both* billing and technical?
- The responder re-sends the full history every call. What happens at turn 200?
- Why prepend the system prompt instead of putting it in `messages` once?

<details>
<summary>Reference solution</summary>

```python
from __future__ import annotations

from typing import Literal

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import SystemMessage
from langgraph.graph import END, START, MessagesState, StateGraph

CATEGORIES = ("billing", "technical", "general")

CLASSIFIER_PROMPT = (
    "Classify the user's support request as exactly one of: "
    "billing, technical, general. Reply with the single word only."
)

PERSONAS = {
    "billing": "You are a billing specialist. Be precise about money.",
    "technical": "You are a support engineer. Be concrete and practical.",
    "general": "You are a friendly support agent.",
}


class SupportState(MessagesState):
    category: str


def make_classifier(model: BaseChatModel):
    def classify(state: SupportState) -> dict:
        reply = model.invoke([SystemMessage(CLASSIFIER_PROMPT), *state["messages"]])
        label = str(reply.content).strip().lower().strip(".!？?\"' ")
        if label not in CATEGORIES:
            label = "general"
        return {"category": label}

    return classify


def make_responder(model: BaseChatModel, category: str):
    def respond(state: SupportState) -> dict:
        reply = model.invoke([SystemMessage(PERSONAS[category]), *state["messages"]])
        return {"messages": [reply]}

    return respond


def route_by_category(state: SupportState) -> Literal["billing", "technical", "general"]:
    return state["category"]


def build_graph(model: BaseChatModel):
    builder = StateGraph(SupportState)

    builder.add_node("classify", make_classifier(model))
    for category in CATEGORIES:
        builder.add_node(category, make_responder(model, category))

    builder.add_edge(START, "classify")
    builder.add_conditional_edges("classify", route_by_category)
    for category in CATEGORIES:
        builder.add_edge(category, END)

    return builder.compile()
```

</details>

---

Next: [05 — Tool calling from scratch](../05-tool-calling) — build the ReAct loop by hand.
