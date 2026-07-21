# 02 — Reducers: how updates merge

> **Concepts:** `Annotated` reducers · `operator.add` · `add_messages` · custom reducers · message ids

## The idea

In exercise 01 the default merge behaviour was invisible because every node
owned a different key. The default is simply:

```python
state[key] = update[key]     # last writer wins
```

That is wrong for anything that should *accumulate* — chat history, logs,
retrieved documents, results from parallel branches. A **reducer** replaces that
assignment with a function of your choosing:

```python
state[key] = reducer(state[key], update[key])
```

You attach it in the type annotation, which means **the merge policy is part of
the state schema**, not scattered across the nodes:

```python
class RunState(TypedDict):
    status: str                                    # no reducer -> overwrite
    log: Annotated[list[str], operator.add]        # list + list -> concatenate
    high_score: Annotated[int, keep_max]           # your own function
```

A reducer is any `(current, update) -> new`. On the very first write `current`
is not `None` but the annotated type's empty value — `[]` for a list, `0` for an
int — which is explained in full below.

### Why `add_messages` and not `operator.add`

`operator.add` on a message list would append blindly. `add_messages` does three
extra things you want for chat history:

1. **Appends** new messages, like `operator.add`.
2. **Replaces by id** — an incoming message whose `id` matches an existing one
   overwrites it instead of duplicating. This is what makes streaming partial
   responses and editing a message work.
3. **Coerces** loose input — `("user", "hi")` or a bare string becomes a proper
   `HumanMessage`.

This pattern is so common that LangGraph ships `MessagesState`, a prebuilt
`TypedDict` that is exactly `{"messages": Annotated[list, add_messages]}`. You
will use it from exercise 04 on; here you build it by hand once so it is not
magic.

### The parallel-write connection

Reducers are also what make **parallel branches** safe. If two nodes run
concurrently and both write `results`, a reducer defines the merge. Without one,
LangGraph raises `InvalidUpdateError` rather than silently letting one branch
clobber the other. Exercise 10 leans on this hard.

### What a reducer actually configures: the channel

Exercise 01 mentioned that each state key becomes a **channel** with a write
policy. A reducer *is* that write policy. Compile a graph and look:

```python
class RunState(TypedDict):
    status: str
    log: Annotated[list[str], operator.add]

>>> graph.channels
{'status': <LastValue>, 'log': <BinaryOperatorAggregate>}
```

Two different objects, chosen by the annotation:

- **`LastValue`** — the default. A write replaces the stored value. It also
  enforces a rule that matters later: *one write per step*. Two nodes writing
  the same `LastValue` key in the same superstep is an `InvalidUpdateError`, not
  a race — LangGraph refuses to pick a winner (exercise 10).
- **`BinaryOperatorAggregate`** — what `Annotated[..., fn]` produces. It holds
  the current value and folds each incoming write into it with your function.

So the annotation is not a hint or a decoration that LangGraph happens to
respect. It **selects the storage strategy** for that key. This is why the merge
policy belongs on the schema rather than in the nodes: it is a property of the
*key*, not of whoever happens to write it, and every node writing `log` gets the
same behaviour whether it knows about the reducer or not.

### The reducer contract, precisely

```python
def keep_max(current: int | None, update: int) -> int:
```

**It is called once per write**, not once at the end. Two nodes writing `log`
means two calls, each folding one update into the running value.

**What is `current` on the very first write?** Not `None`, despite what the
`int | None` in `keep_max`'s signature suggests. `BinaryOperatorAggregate` seeds
the channel by calling your annotated type with no arguments:

| Annotation | `current` on first write |
|---|---|
| `Annotated[list[str], operator.add]` | `[]` |
| `Annotated[int, keep_max]` | `0` |
| `Annotated[str, operator.add]` | `""` |

Verify it yourself — put a `print` in `keep_max` and invoke the graph without
seeding `high_score`. You will see `0`, then `10`.

Two consequences. First, the annotated type must be **constructible with no
arguments**, which is why you annotate `list[str]` and not something exotic.
Second, inside a reducer you **cannot distinguish "never written" from "written
the zero value"** — if that difference matters to your logic, the difference
belongs in a separate key, not in a sentinel.

So why does the spec still test `keep_max(None, 4) == 4`? Because a reducer is
an ordinary function with an ordinary contract. It gets called directly in unit
tests, it can be reused on a differently-typed key, and defending one `if` is
cheaper than depending on a framework's seeding behaviour. Handle `None`; just
do not believe it is the normal case.

**It must not mutate `current`.** Return a new value. `current.append(x)` edits
an object the framework may still be holding elsewhere — most visibly in
checkpoints (exercise 07), where the "previous" state suddenly contains the new
item too. `operator.add` on lists is safe precisely because `+` builds a new
list.

**It should be associative.** For a linear graph, updates fold in a predictable
order and anything works. Once branches run in parallel (exercise 10), the
grouping of merges is up to the scheduler, so `(a + b) + c` must equal
`a + (b + c)`. Concatenation and `max` satisfy this; "subtract" or "divide" do
not, and a reducer like that will produce results that change between runs.

Commutativity — whether `a + b == b + a` — is a subtler question. List
concatenation is *not* commutative, which is exactly why `log` and `messages`
preserve order. Relying on that order across parallel branches is where it gets
risky; see exercise 10's *Think about it*.

### `add_messages` in depth

It is a reducer like any other, but with three behaviours tuned for
conversation. Worth knowing individually, because each one solves a specific
problem you would otherwise hit:

**1. It appends.** The baseline, same as `operator.add`.

**2. It upserts by `id`.** Every LangChain message has an `.id`. If an incoming
message carries an id that already exists in the list, `add_messages`
**replaces** that message instead of appending a copy. This is what makes two
things work:

- *Streaming a reply.* Partial versions of the same message arrive repeatedly
  and collapse onto one entry rather than producing forty duplicates.
- *Editing history.* Re-emit a message with an existing id and you have
  rewritten that turn — the basis of "edit your message and regenerate".

It also produces a genuinely confusing bug the first time you meet it, and you
will meet it in exercise 05: replay the *same* `AIMessage` object twice and the
second one overwrites the first, because it has the same id. The list does not
grow, and a loop that looks infinite silently terminates. When message counts
do not match what you expect, check ids first.

**3. It coerces.** `("user", "hi")` becomes a `HumanMessage`; a bare string
becomes one too. That is why every test in this course can write
`{"messages": [("user", "hello")]}` instead of importing a message class. The
coercion happens in the reducer, so it applies to node returns as well as
invoke inputs.

There is a fourth behaviour you will not need here but should know exists:
returning a `RemoveMessage(id=...)` **deletes** that message from the list. That
is how trimming and summarising are implemented — a node that returns removals
for old turns, which is the standard answer to the context-growth problem raised
in exercises 04 and 07.

### The bugs this design prevents — and the ones it invites

| Symptom | Cause |
|---|---|
| `TypeError: can only concatenate list (not "str") to list` | A node returned `"entered one"` where the channel holds a list. The reducer literally ran `[...] + "entered one"`. **Return the container type the channel holds**, always — `["entered one"]`. |
| History has duplicates | Two distinct message objects with the same content. Ids are per-object, so distinct objects are distinct entries. |
| History is missing entries | The opposite: the *same* message object written twice, so `add_messages` upserted over the earlier copy. Bites in exercise 05. |
| `InvalidUpdateError: can receive only one value per step` | Two nodes wrote a no-reducer (`LastValue`) key in one superstep. Exercise 10. |
| A reducer never runs | Nothing wrote that key. Reducers fire on writes, not on reads or on every step. |

That first one is worth dwelling on, because the failure is *loud*. A framework
that silently splatted the string into `['e','n','t',...]` would be far worse:
you would find out days later, in a log nobody reads. An immediate `TypeError`
naming both types is the good outcome.

### `MessagesState`, and why you are building it by hand first

`{"messages": Annotated[list, add_messages]}` is so ubiquitous that LangGraph
ships it prebuilt:

```python
from langgraph.graph import MessagesState

class SupportState(MessagesState):   # exercise 04 onward
    category: str
```

From exercise 04 you will use it constantly and never think about it again.
Building it once by hand here is the point: when a message mysteriously fails to
appear in `state["messages"]` three exercises from now, you will know there is a
reducer doing id-matching underneath, rather than treating it as framework
weather.

## The imports you need

```python
from __future__ import annotations

import operator
from typing import Annotated, TypedDict

from langchain_core.messages import AIMessage, HumanMessage
from langgraph.graph import END, START, StateGraph, add_messages
```

Carrying over from 01: `from __future__ import annotations`, `TypedDict`,
`StateGraph`, `START`, `END`. New in this exercise:

**`operator`** — the standard-library module exposing Python's operators as
functions. `operator.add(a, b)` is exactly `a + b`, and for two lists `+` means
concatenate. So `operator.add` *is* a valid reducer, with no wrapper needed —
you are just handing LangGraph the `+` function.

**`Annotated`** — the typing construct that attaches **metadata** to a type:

```python
Annotated[list[str], operator.add]
#         ^ the type  ^ the metadata
```

To Python, `Annotated[X, anything]` still just means `X` — the second argument
is ignored by the type system entirely. It is a place to hang extra information
that *some other library* will read. LangGraph reads it and finds your reducer.
That is the whole trick: the reducer travels with the type, so the state schema
declares both what a key holds and how it merges.

**`add_messages`** — LangGraph's purpose-built reducer for message lists
(appending, id-replacement, and coercion, as above). Note where it comes from:
`langgraph.graph`, alongside `StateGraph`.

**`AIMessage` / `HumanMessage`** — LangChain's message classes. A conversation
is a list of typed messages rather than strings, because a model needs to know
*who said what*: `HumanMessage` is the user, `AIMessage` is the model,
`SystemMessage` is instruction (exercise 04), `ToolMessage` is a tool result
(exercise 05). Here you only need `AIMessage` for the nodes; `HumanMessage` is
for the optional demo's seed state.

They come from `langchain_core`, not `langgraph` — worth noticing. LangGraph
owns *orchestration*; LangChain Core owns *messages, models, and tools*. Two
packages, two jobs.

## What to build

### 1. `keep_max(current, update)`

A custom reducer returning whichever value is larger. Signature:

```python
def keep_max(current: int | None, update: int) -> int:
```

Handle `current is None` by returning `update`; otherwise return
`max(current, update)`. In *this* graph LangGraph will actually hand you `0`
rather than `None` on the first write (see above), but the spec unit-tests the
`None` case directly, and a reducer that defends its own contract is the right
habit.

### 2. `RunState`

A `TypedDict` with four keys, each with a different merge behaviour:

| Key | Annotation | Behaviour |
|---|---|---|
| `status` | `str` | no reducer — last writer wins |
| `log` | `Annotated[list[str], operator.add]` | every node appends |
| `messages` | `Annotated[list, add_messages]` | chat history semantics |
| `high_score` | `Annotated[int, keep_max]` | your own reducer |

### 3. `stage_one` and `stage_two`

Two nodes, each returning **all four keys**.

`stage_one`:

| Key | Value |
|---|---|
| `status` | `"one"` |
| `log` | `["entered one"]` |
| `messages` | `[AIMessage("hello from one")]` |
| `high_score` | `10` |

`stage_two`:

| Key | Value |
|---|---|
| `status` | `"two"` |
| `log` | `["entered two"]` |
| `messages` | `[AIMessage("hello from two")]` |
| `high_score` | `5` — lower on purpose, so `keep_max` has something to ignore |

Watch for the classic mistake: with `operator.add` a node must return a **list**
(`{"log": ["entered one"]}`), not a bare string. Returning `"entered one"` makes
the reducer evaluate `[...] + "entered one"`, which raises
`TypeError: can only concatenate list (not "str") to list`.

### 4. `build_graph()`

`START → one → two → END`, with nodes named `one` and `two`. Compile and return.

### 5. Optional: a runnable demo

```python
if __name__ == "__main__":
    from labgraph import print_state

    graph = build_graph()
    result = graph.invoke(
        {
            "status": "start",
            "log": ["initial"],
            "messages": [HumanMessage("hi")],
            "high_score": 0,
        }
    )
    print_state(result, title="after two stages")
    print("\nNotice: status was overwritten, log/messages accumulated,")
    print("and high_score kept 10 even though stage_two wrote 5.")
```

## Run it

```bash
uv run pytest exercises/02-reducers -v
uv run python exercises/02-reducers/graph.py
```

## Think about it

- Why is the reducer on the *state schema* rather than on the node that writes it?
- `keep_max` gets `None` on first write. Why not just seed the key with `0`?
- If two parallel nodes both write a key with no reducer, what should happen —
  and what does LangGraph actually do?
- Can a reducer be non-commutative? Should it be, given parallel branches can
  merge in any order?
- `Annotated[list, add_messages]` — what would break if you wrote
  `Annotated[list, add_messages()]` instead?

<details>
<summary>Reference solution</summary>

```python
from __future__ import annotations

import operator
from typing import Annotated, TypedDict

from langchain_core.messages import AIMessage
from langgraph.graph import END, START, StateGraph, add_messages


def keep_max(current: int | None, update: int) -> int:
    if current is None:
        return update
    return max(current, update)


class RunState(TypedDict):
    status: str
    log: Annotated[list[str], operator.add]
    messages: Annotated[list, add_messages]
    high_score: Annotated[int, keep_max]


def stage_one(state: RunState) -> dict:
    return {
        "status": "one",
        "log": ["entered one"],
        "messages": [AIMessage("hello from one")],
        "high_score": 10,
    }


def stage_two(state: RunState) -> dict:
    return {
        "status": "two",
        "log": ["entered two"],
        "messages": [AIMessage("hello from two")],
        "high_score": 5,
    }


def build_graph():
    builder = StateGraph(RunState)
    builder.add_node("one", stage_one)
    builder.add_node("two", stage_two)
    builder.add_edge(START, "one")
    builder.add_edge("one", "two")
    builder.add_edge("two", END)
    return builder.compile()
```

</details>

---

Next: [03 — Conditional routing](../03-conditional-routing) — graphs that branch and loop.
