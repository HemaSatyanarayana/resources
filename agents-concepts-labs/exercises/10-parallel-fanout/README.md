# 10 — Parallel fan-out: `Send` and map-reduce

> **Concepts:** `Send` · dynamic parallelism · worker state schemas · supersteps · reducers under concurrency · `InvalidUpdateError`

## The idea

Every graph so far had a shape you could draw before running it. This one does
not:

```python
return [Send("research", {"subtopic": s}) for s in state["subtopics"]]
```

Three subtopics, three concurrent copies of `research`. Ten, ten. Zero, none.
The planner decides at runtime, and the graph widens to match.

```
START -> plan ──┬─> research("history") ─┐
                ├─> research("ethics")  ─┼─> write -> END
                └─> research("jobs")    ─┘
```

That is map-reduce, and both halves are things you already have:

| | This exercise | From |
|---|---|---|
| **map** | `Send` — N tasks, one per item | new |
| **reduce** | `Annotated[list[str], operator.add]` | exercise 02 |

### A `Send` is a task, not an edge

An edge says *"go to this node"*. A `Send` says *"run this node, with this
input"* — and you can issue as many as you like to the same node. Returning
three Sends for `"research"` does not visit one node three times in sequence;
it schedules three independent tasks in a single superstep.

### The payload replaces the state

```python
Send("research", {"subtopic": "ethics"})
```

`{"subtopic": "ethics"}` is the worker's **entire input**. It cannot see
`topic`, it cannot see `notes`, it cannot see what its siblings are doing. Hence
the separate `WorkerState` schema — a worker that reads a different shape than
the parent is a feature, not a mistake.

That isolation is what makes concurrency safe here. Workers cannot race on
something they cannot reach. If a worker needs the topic too, you must put it in
the payload — explicitly, per worker.

### Concurrent writes need a reducer — or it is an error

Three workers return `{"notes": [...]}` in the same superstep. With
`Annotated[list[str], operator.add]` they concatenate. Without it:

```
InvalidUpdateError: At key 'notes': Can receive only one value per step.
```

LangGraph will not guess which write wins. This is exercise 02's lesson with the
training wheels off: a reducer is not a convenience for accumulating history, it
is **the merge strategy for concurrent writes**, and fan-out makes it mandatory.
The spec proves both halves with two self-contained graphs.

### One superstep, then a join

```python
builder.add_edge("research", "write")
```

One ordinary edge, even though `research` runs N times. LangGraph waits for
**all** N to finish, then runs `write` **once**. You do not write barrier or
join logic; the superstep model is the join.

### The empty case will bite you

An empty list of Sends advances nothing — the graph has nowhere to go. So handle
it explicitly:

```python
if not subtopics:
    return "write"        # a plain node name is also a legal return
```

A conditional edge may return Sends *or* a node name. And because LangGraph
cannot infer the targets of a `Send` from a type annotation, spell them out:

```python
builder.add_conditional_edges("plan", fan_out, ["research", "write"])
```

Leave that list off and your graph still runs, but `get_graph()` cannot draw the
branch — the diagram quietly lies about what your agent does.

### Static topology, dynamic width

Every graph so far was fully described before it ran. `add_node` and `add_edge`
happen at build time, once, at import — so "how many nodes are there?" always had
an answer you could read off the source.

`Send` breaks that, but only in one specific dimension. The **topology** is still
static: there is exactly one node named `research`, registered once. What becomes
dynamic is **how many tasks run against it**:

```python
builder.add_node("research", make_researcher(model))   # ONE node, at build time
...
[Send("research", {...}) for s in subtopics]           # N tasks, at run time
```

There is no such thing as "adding N nodes". Three concurrent researchers are
three **tasks** dispatched to one node definition, in one superstep — the same
relationship as one function and three concurrent calls to it.

That distinction explains the whole API. It is why `add_conditional_edges` needs
`["research", "write"]` spelled out (the *possible* targets are static and
knowable), and why the worker takes a `WorkerState` (each task gets its own
input, like arguments to a call).

### Why a `Send` payload replaces the state entirely

```python
Send("research", {"subtopic": "ethics"})
```

The worker receives `{"subtopic": "ethics"}` and nothing else. Not `topic`, not
`notes`, not what its siblings are doing.

This looks restrictive until you consider the alternative. If workers could read
the full parent state, then whether worker 2 saw worker 1's note would depend on
which finished first — a data race, non-deterministic, and untestable. By
handing each task an isolated input, LangGraph makes that class of bug
unrepresentable: **workers cannot race on state they cannot reach.**

The price is explicitness. A worker that needs the parent topic must be *given*
it, per task:

```python
Send("research", {"subtopic": s, "topic": state["topic"]})
```

That is a feature. Anything a worker depends on appears at the dispatch site,
where you can see it, rather than being reached for from a shared scope.

Hence the second schema. `WorkerState` is not a subset of `ResearchState` and is
not derived from it — it is the **argument type of one task**. A node whose input
schema differs from the graph's state schema is a normal thing here, not a
mistake.

### Concurrent writes: where exercise 02 becomes non-negotiable

Three workers finish in the same superstep and all three return `{"notes": [...]}`.
Three writes to one channel, one step.

With `Annotated[list[str], operator.add]`, the channel is a
`BinaryOperatorAggregate` (exercise 02): it folds each write into the running
value, and all three notes survive. Without the annotation it is a `LastValue`
channel, which permits exactly one write per step:

```
InvalidUpdateError: At key 'notes': Can receive only one value per step.
```

Note what LangGraph does **not** do: it does not pick the last writer, or the
first, or warn and continue. Silently keeping one of three results would be a
data-loss bug that surfaces as "the report is missing sections, sometimes" —
the worst kind. Refusing to guess is the right call, and it is why the spec
proves both halves with two self-contained graphs.

So: in exercise 02 a reducer was a convenience for accumulating a log. Here it
is **the merge strategy for concurrent writes**, and fan-out makes it mandatory.
Same syntax, entirely different stakes — and this is also where a reducer's
associativity (exercise 02) stops being theoretical, because the order in which
those three writes fold together is the scheduler's business, not yours.

### The join is free

```python
builder.add_edge("research", "write")
```

One ordinary edge, and yet `write` runs **once**, after **all** N workers finish
— not once per worker.

You wrote no barrier, no counter, no "have they all finished yet?" check. That
falls out of the superstep model (exercise 01): a superstep completes when every
task scheduled in it has completed, *then* the next is scheduled. Synchronisation
is the execution model, not something you implement.

Watch the update stream and the structure is visible:

```
{"plan":     {...}}      superstep 1
{"research": {...}}   ┐
{"research": {...}}   ├ superstep 2 — all workers, one step
{"research": {...}}   ┘
{"write":    {...}}      superstep 3 — runs once
```

Each worker emits its **own** update chunk (exercise 09) — they are separate
tasks — but they belong to one superstep, and `write` waits for the whole
superstep.

### What this does not handle

The exercise is a clean map-reduce. Production fan-out is where the interesting
failures live, and it is worth knowing what you have *not* solved:

**Unbounded width.** `Send` launches as many tasks as the planner names. Fifty
subtopics is fifty concurrent model calls, which will hit provider rate limits
and may cost more in one run than you expected to spend in a day. Real systems
cap the width — chunk the list, or cap it at the planner.

**Failure semantics.** One worker raising does not leave the other nine's work
neatly available; the superstep fails. Whether you want all-or-nothing or
best-effort is a decision, and best-effort means each worker catching its own
exceptions and returning a note describing the failure — the "errors as
observations" habit from exercise 05, applied to workers.

**Ordering.** Notes come back in `Send` order here, and the tests rely on it.
Be careful about depending on that in general: it is a property of how tasks are
dispatched, not a documented guarantee, and a retried worker will not politely
hold its place. If order matters to your output, sort explicitly on something
you control.

**Isolation cuts both ways.** Workers cannot see each other — which is what makes
this safe, and also what makes it the wrong shape when research task 2 genuinely
needs what task 1 found. That is a sequential dependency, and forcing it into a
fan-out produces workers that duplicate each other's effort. Fan out over
*independent* work; chain dependent work.

**The reduce step is naive.** `write_report` concatenates. It does not
deduplicate, reconcile contradictions between workers, or notice that two
subtopics returned the same fact. A real reducer often needs its own model call
— at which point "reduce" is just another LLM node, and you are back to
everything from exercise 04.

## The imports you need

```python
from __future__ import annotations

import operator
from typing import Annotated, TypedDict

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import SystemMessage
from langgraph.graph import END, START, StateGraph
from langgraph.types import Send
```

Almost all of this is exercise 02 and 04 material returning. Two things worth
noticing:

**`Send`** — from `langgraph.types`, the same module `interrupt` and `Command`
came from in 08. All three are *runtime control-flow* objects rather than
building blocks, which is why they share a home. A `Send` is a pair:

```python
Send("research", {"subtopic": "ethics"})
#     ^ node name  ^ that node's entire input
```

You return them from a conditional edge instead of a node name. Returning three
schedules three tasks.

**`operator` and `Annotated` are back — and this time they are mandatory.** In
02 the reducer on `log` was a convenience for accumulating history. Here,
several workers write `notes` in the *same step*, and without a reducer
LangGraph raises `InvalidUpdateError` rather than guessing which write wins.
Same syntax, much sharper consequence.

**No `MessagesState`.** This graph's state is not a conversation — it is a topic,
a plan, some notes, and a report. Reach for a plain `TypedDict` when messages are
not the point; `MessagesState` is a convenience for chat-shaped state, not a
requirement.

## What to build

### 1. Two prompts

```python
PLANNER_PROMPT = (
    "List 3 short subtopics for researching the user's topic. "
    "Reply with a comma-separated list and nothing else."
)

RESEARCHER_PROMPT = "Write one short factual sentence about the subtopic given."
```

### 2. Two state schemas

`ResearchState` — the parent:

| Key | Type |
|---|---|
| `topic` | `str` |
| `subtopics` | `list[str]` |
| `notes` | `Annotated[list[str], operator.add]` ← **the line that makes fan-out legal** |
| `report` | `str` |

`WorkerState` — what **one worker** sees, and it is not `ResearchState`:

```python
class WorkerState(TypedDict):
    subtopic: str
```

A `Send` payload becomes the node's entire input. The worker cannot see `topic`,
`notes`, or its siblings' work. That isolation is the point — it is what makes
running them at once safe.

### 3. `parse_subtopics(raw) -> list[str]`

Split a comma-separated model reply into clean subtopics: strip whitespace, drop
empties. `"a,,b"` → `["a", "b"]`, `"   "` → `[]`. Same defensive parsing as
exercise 04 — the model's output is never as tidy as the prompt asked for.

### 4. `make_planner(model)`

A node factory. The node calls the model with `SystemMessage(PLANNER_PROMPT)`
and the topic as the user message, then returns
`{"subtopics": parse_subtopics(reply)}`.

### 5. `fan_out(state)`

The conditional edge that decides how many workers to launch:

```python
return [Send("research", {"subtopic": s}) for s in state["subtopics"]]
```

When there are **no** subtopics, return the plain string `"write"` instead. A
conditional edge may return Sends *or* a node name, and an empty list of Sends
is a dead end with nothing to advance the graph.

### 6. `make_researcher(model)`

The worker factory. Note the inner signature takes a `WorkerState`, not the
parent. Call the model with `SystemMessage(RESEARCHER_PROMPT)` and the subtopic,
then return:

```python
{"notes": [f"{subtopic}: {reply.content}"]}
```

A **list of one**. The reducer concatenates lists; hand it a bare string and it
raises `TypeError: can only concatenate list (not "str") to list` (exercise 02).

### 7. `write_report(state)`

The reduce step: every worker has finished and `notes` holds all of them. Return
`{"report": ...}` shaped like:

```
# <topic>

- <note>
- <note>
```

With no notes, the body is `(no findings)`.

### 8. `build_graph(model)`

```
START    -> plan
plan     -> (conditional: fan_out) -> N x research | write
research -> write
write    -> END
```

Node names: `plan`, `research`, `write`. One wrinkle — the conditional edge
needs its possible targets spelled out, because `fan_out` returns `Send` objects
that LangGraph cannot infer from a type annotation:

```python
builder.add_conditional_edges("plan", fan_out, ["research", "write"])
```

And note `research -> write` is a **single ordinary edge**, even though
`research` runs N times. LangGraph waits for all N before `write` starts.

### 9. Optional: a runnable demo

```python
if __name__ == "__main__":
    from labgraph import get_chat_model

    graph = build_graph(get_chat_model())
    result = graph.invoke(
        {"topic": "the Apollo program", "subtopics": [], "notes": [], "report": ""}
    )

    print(f"planned: {result['subtopics']}\n")
    print(result["report"])
```

The test worth reading twice is `test_worker_sees_its_subtopic_and_nothing_else`
— it asserts the parent topic is *absent* from the worker's prompt. Worker
isolation is a property you should be able to state precisely.

## Run it

```bash
uv run pytest exercises/10-parallel-fanout -v
uv run python exercises/10-parallel-fanout/graph.py     # needs .env
```

## Think about it

- All N workers hit your model provider at once. What does subtopic #47 do to
  your rate limit, and where would you cap the width?
- One worker raises. What happens to the other nine, and to `write`? What
  *should* happen?
- The notes come back in Send order here. Would you rely on that? What if a
  worker retried?
- Workers cannot see each other. When is that wrong — when does research task 2
  genuinely need what task 1 found?
- Each worker is a node. What if each worker were a whole *agent* — a subgraph
  with its own tools and loop? (Nothing stops you: `add_node` takes a compiled
  graph.)
- `write` blindly concatenates. What would a real reduce step do about
  duplicated or contradictory notes?

<details>
<summary>Reference solution</summary>

```python
from __future__ import annotations

import operator
from typing import Annotated, TypedDict

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import SystemMessage
from langgraph.graph import END, START, StateGraph
from langgraph.types import Send

PLANNER_PROMPT = (
    "List 3 short subtopics for researching the user's topic. "
    "Reply with a comma-separated list and nothing else."
)

RESEARCHER_PROMPT = "Write one short factual sentence about the subtopic given."


class ResearchState(TypedDict):
    topic: str
    subtopics: list[str]
    notes: Annotated[list[str], operator.add]
    report: str


class WorkerState(TypedDict):
    subtopic: str


def parse_subtopics(raw: str) -> list[str]:
    return [part.strip() for part in raw.split(",") if part.strip()]


def make_planner(model: BaseChatModel):
    def plan(state: ResearchState) -> dict:
        reply = model.invoke([SystemMessage(PLANNER_PROMPT), ("user", state["topic"])])
        return {"subtopics": parse_subtopics(str(reply.content))}

    return plan


def fan_out(state: ResearchState):
    subtopics = state.get("subtopics") or []
    if not subtopics:
        return "write"
    return [Send("research", {"subtopic": subtopic}) for subtopic in subtopics]


def make_researcher(model: BaseChatModel):
    def research(state: WorkerState) -> dict:
        reply = model.invoke([SystemMessage(RESEARCHER_PROMPT), ("user", state["subtopic"])])
        return {"notes": [f"{state['subtopic']}: {reply.content}"]}

    return research


def write_report(state: ResearchState) -> dict:
    notes = state.get("notes") or []
    if not notes:
        return {"report": f"# {state['topic']}\n\n(no findings)"}
    body = "\n".join(f"- {note}" for note in notes)
    return {"report": f"# {state['topic']}\n\n{body}"}


def build_graph(model: BaseChatModel):
    builder = StateGraph(ResearchState)

    builder.add_node("plan", make_planner(model))
    builder.add_node("research", make_researcher(model))
    builder.add_node("write", write_report)

    builder.add_edge(START, "plan")
    builder.add_conditional_edges("plan", fan_out, ["research", "write"])
    builder.add_edge("research", "write")
    builder.add_edge("write", END)

    return builder.compile()
```

`make_researcher` is registered **once** as a single node. The three concurrent
copies are three tasks against one node definition — there is no such thing as
"adding N nodes" here.

</details>

---

That is the course. See the [top-level README](../../README.md) for where to go
next — subgraphs, multi-agent supervisors, long-term memory stores, and
deploying a graph.
