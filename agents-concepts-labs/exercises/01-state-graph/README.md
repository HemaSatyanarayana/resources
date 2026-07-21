# 01 — Your first StateGraph

> **Concepts:** state schema · nodes · edges · `START`/`END` · compile · partial updates

## The idea

Most agent frameworks hide control flow inside a loop you cannot see. LangGraph
does the opposite: you declare the flow as a **graph**, and the framework runs
it. Three pieces, and nothing else:

| Piece | What it is |
|---|---|
| **State** | A schema (usually a `TypedDict`) describing everything the run knows. |
| **Node** | A plain function `state -> partial update`. |
| **Edge** | "After this node, go to that one." |

The execution model in one paragraph: you `invoke` the graph with an initial
state. LangGraph runs the node wired to `START`, takes the dict that node
returned, and **merges** it into the state. Then it follows that node's outgoing
edge, runs the next one, merges again — until it reaches `END`, at which point
the merged state is handed back to you.

Two consequences worth internalising now, because everything later builds on them:

**Nodes never call each other.** They communicate only through state. `count`
does not invoke `clean`; it just reads the `cleaned` key and trusts that whoever
ran before it put a value there. That is what makes nodes independently
testable, retryable, and resumable.

**Nodes return partial updates.** A node returns `{"cleaned": "..."}`, not the
whole state. Anything you leave out is untouched. By default a returned key
*overwrites* the existing value — how to make it accumulate instead is exercise
02.

```
START ──▶ clean ──▶ count ──▶ summarize ──▶ END
          │          │          │
          └──────────┴──────────┴──▶ all reading/writing one shared state dict
```

### Build time vs. run time

Two phases, and confusing them is the first mistake everyone makes:

```python
builder = StateGraph(PipelineState)   # build time: describe the graph
builder.add_node("clean", clean_text) # still build time
graph = builder.compile()             # freeze it — now it is runnable
graph.invoke({...})                   # run time: data flows through
```

`add_node` does not run your function. It **registers** it. Nothing executes
until `invoke`. A `StateGraph` is a blueprint; `compile()` turns it into the
thing that runs, and that compiled object is what you call, stream, and — from
exercise 07 — attach memory to.

### What actually happens when you call `invoke`

"LangGraph runs the nodes" is too vague to debug with. Here is the real
sequence, for `graph.invoke({"text": "  Hello   WORLD  "})`:

**Step 0 — seed.** Your input dict is written into the state. Keys you did not
supply simply have no value yet; they are not initialised to `""` or `0`. This
is why a node that reads `state["cleaned"]` before anything wrote it raises
`KeyError` rather than quietly seeing an empty string.

**Step 1 — find the entry.** LangGraph looks at what `START` points to. That is
`clean`, so `clean` is scheduled.

**Step 2 — run the node.** `clean_text(state)` is called with the *whole* state
as a plain dict. It returns `{"cleaned": "hello world"}`.

**Step 3 — merge.** LangGraph takes that returned dict and applies each key to
the state. With no reducer configured (exercise 02), applying means assignment.
State is now `{"text": "...", "cleaned": "hello world"}`.

**Step 4 — follow the edge.** `clean`'s outgoing edge points to `count`.
Schedule it. Repeat steps 2–4.

**Step 5 — stop.** When the scheduled node is `END`, the run finishes and the
accumulated state is returned to you — *all* of it, including `text`, which no
node ever wrote back.

Each pass through steps 2–4 is called a **superstep**. Right now every superstep
runs exactly one node, so the distinction looks pedantic. It stops being
pedantic in exercise 10, where one superstep runs fifty nodes at once, and in
exercise 03, where the framework's runaway-loop protection counts *supersteps*
rather than node calls.

### Under the hood: state is a set of channels

The mental model of "a dict LangGraph passes around" is close enough to get
started, but the real implementation is worth knowing because it explains
exercise 02 completely.

Each key in your schema becomes a **channel** — an independent slot with its own
write policy. You can see them on a compiled graph:

```python
>>> graph.channels
{'text': <LastValue>, 'cleaned': <LastValue>, 'word_count': <LastValue>, ...}
```

`LastValue` is the default policy, and it means exactly what it says: a write
replaces whatever was there. A node returning `{"cleaned": x}` is not mutating a
shared dictionary — it is **writing to the `cleaned` channel**, and the channel
decides what that write means.

That indirection is the whole design. Change a key's channel and you change how
writes to it behave, without touching a single node. Exercise 02 does precisely
that.

### `compile()` is a validation step, not a formality

```python
graph = builder.compile()
```

`compile()` walks the graph you described and rejects it if it cannot run —
an edge pointing at a node name that does not exist, or a node that nothing can
ever reach. Catching that at build time, once at import, beats discovering it
on a production request.

What comes back is a `CompiledStateGraph`, and its ancestry is informative:

```
CompiledStateGraph -> Pregel -> Runnable
```

**`Pregel`** is Google's graph-processing model — the source of the superstep
vocabulary above. LangGraph is an implementation of it aimed at agents.

**`Runnable`** is LangChain's universal interface, and it is why the compiled
object supports far more than `.invoke()`: `.stream()` (exercise 09),
`.batch()`, and async twins `.ainvoke()` / `.astream()` come with the type. You
get all of them without doing anything.

A `StateGraph` is mutable and describes intent; a `CompiledStateGraph` is frozen
and executes. Build once at import; invoke many times per request.

### The node contract, precisely

A node is any callable matching `(state) -> dict | None`. Beyond that:

**It receives the whole state**, not the keys it declared interest in — there is
no such declaration. Read what you like.

**It returns a partial update**, and the keys must exist in your schema.
Returning `{"typo_key": 1}` is an error, not a silently ignored write. Your
schema is the contract.

**Returning `None` (or `{}`) is legal** and means "I changed nothing". A node
that only logs, validates, or emits a progress event is a perfectly good node.

**Do not mutate the state dict in place.** `state["cleaned"] = x` may appear to
work today, but it bypasses the channel machinery — which means no reducer
runs, no checkpoint records the change (exercise 07), and nothing streams it
(exercise 09). Return the update; let the framework apply it.

**Assume it can be re-run.** Not in this exercise, but soon: a retry, a resumed
run (07), or a node containing an interrupt (08) will each execute the same node
twice. Nodes that only compute and return survive that. Nodes that also charge a
credit card do not.

### Why bother with a graph at all?

For this pipeline, three chained function calls would be shorter, and it is
worth being honest about what the ceremony buys — because for straight-line code
the answer is "not much".

It pays off the moment you want any of the following, all of which are
effectively free once the flow is *data* rather than *control flow*:

| You want | Why the graph gives it to you |
|---|---|
| **Resume a run** | State is external and serialisable, so a run can stop and continue later (07) |
| **A human in the middle** | Pausing is just "don't schedule the next node yet" (08) |
| **Progress and observability** | The framework knows which node is running; it can stream that (09) |
| **Parallelism** | Independent nodes are visibly independent, so they can be scheduled together (10) |
| **Cycles with a budget** | "Go back to a previous node" is an edge, not a `while` loop you must reason about (03) |

None of that is available to a plain function pipeline without you building it
by hand. That is the trade: some indirection now, in exchange for capabilities
that are otherwise very expensive to retrofit.

## The imports you need

Open your empty `graph.py` and start with these. Nothing here is magic — each
line earns its place:

```python
from __future__ import annotations

from typing import TypedDict

from langgraph.graph import END, START, StateGraph
```

**`from __future__ import annotations`** — makes Python treat every type
annotation as a *string* instead of evaluating it on the spot. That lets you
write modern syntax like `list[str]` or `int | None` regardless of Python
version, and it costs nothing at run time. Every file in this course starts with
it. It must be the first import; Python enforces that for `__future__`.

**`TypedDict`** — a dictionary whose keys and value types you declare up front:

```python
class PipelineState(TypedDict):
    text: str
```

At run time it is an ordinary `dict` — nothing is validated, nothing is
enforced. Its job is to *describe* the shape: your editor uses it for
completion, and **LangGraph reads it to learn what keys the state has**. That is
why the schema is a required argument to `StateGraph`.

**`StateGraph`** — the builder you add nodes and edges to. You construct it with
your state schema, so it knows what keys exist.

**`START` and `END`** — two sentinel values (really the strings `"__start__"`
and `"__end__"`) naming the graph's entry and exit. They are *not* nodes you
write; they are markers you connect real nodes to. `add_edge(START, "clean")`
means "begin here"; `add_edge("summarize", END)` means "after this, stop".

For the optional runnable demo at the bottom of the file:

```python
from labgraph import print_state
```

`labgraph/` is the course's helper package — already written, not an exercise.
`print_state` renders a state dict readably.

## What to build

### 1. `PipelineState`

A `TypedDict` with four keys:

| Key | Type | Meaning |
|---|---|---|
| `text` | `str` | the raw input |
| `cleaned` | `str` | normalised text |
| `word_count` | `int` | words in `cleaned` |
| `summary` | `str` | a one-line human-readable report |

### 2. Three node functions

Each takes the whole state and returns **only the key it changed**.

- **`clean_text(state) -> dict`** — lowercase `state["text"]` and collapse
  whitespace, returning `{"cleaned": ...}`.
  `"  Hello   WORLD  "` becomes `"hello world"`.

- **`count_words(state) -> dict`** — return `{"word_count": <int>}`, counting
  words in `state["cleaned"]`. Note it reads a key a *previous* node wrote —
  that is the entire contract between nodes.

- **`summarize(state) -> dict`** — return `{"summary": ...}` in exactly this
  format: `"<word_count> words: <cleaned>"`, e.g. `"2 words: hello world"`.

Returning the whole state would also work, but partial updates are the idiom:
they make it obvious which node owns which key.

### 3. `build_pipeline()`

Wire and compile:

1. `builder = StateGraph(PipelineState)`
2. `builder.add_node("clean", clean_text)`, and likewise `"count"` and
   `"summarize"`. The node **name** is what shows up in streams, checkpoints and
   debug output; it need not match the function name, but the spec pins these
   three, so spell them exactly.
3. Connect with `builder.add_edge(from, to)`:
   `START → clean → count → summarize → END`.
4. `return builder.compile()`

### 4. Optional: a runnable demo

So that `uv run python exercises/01-state-graph/graph.py` shows something:

```python
if __name__ == "__main__":
    from labgraph import print_state

    graph = build_pipeline()
    result = graph.invoke({"text": "  LangGraph makes agent  control flow EXPLICIT  "})
    print_state(result, title="final state")
```

## Run it

```bash
uv run pytest exercises/01-state-graph -v          # the spec
uv run python exercises/01-state-graph/graph.py    # watch it run
```

While `graph.py` is still empty, pytest reports a **collection error** rather
than test failures — there is nothing to import yet. That is the correct
starting point, not something to fix.

## Think about it

- `invoke()` returned `text` even though no node wrote it. Why?
- What happens if two keys are written by one node? If one node writes a key
  another node also writes?
- `clean_text` is a pure function with no LangGraph imports in its body. What
  does that buy you when this pipeline grows to 30 nodes?
- Why does `stream()` yield `{"clean": {...}}` rather than the full state?
- `TypedDict` is not enforced at run time. What stops a node returning
  `{"cleaned": 42}` — and what happens downstream if it does?

<details>
<summary>Reference solution — try first, seriously</summary>

```python
from __future__ import annotations

from typing import TypedDict

from langgraph.graph import END, START, StateGraph


class PipelineState(TypedDict):
    text: str
    cleaned: str
    word_count: int
    summary: str


def clean_text(state: PipelineState) -> dict:
    return {"cleaned": " ".join(state["text"].split()).lower()}


def count_words(state: PipelineState) -> dict:
    return {"word_count": len(state["cleaned"].split())}


def summarize(state: PipelineState) -> dict:
    return {"summary": f"{state['word_count']} words: {state['cleaned']}"}


def build_pipeline():
    builder = StateGraph(PipelineState)

    builder.add_node("clean", clean_text)
    builder.add_node("count", count_words)
    builder.add_node("summarize", summarize)

    builder.add_edge(START, "clean")
    builder.add_edge("clean", "count")
    builder.add_edge("count", "summarize")
    builder.add_edge("summarize", END)

    return builder.compile()
```

`" ".join(text.split())` is the idiomatic whitespace collapse: `split()` with no
argument splits on *runs* of whitespace and drops empties, so leading, trailing,
and doubled spaces all disappear in one step.

</details>

---

Next: [02 — Reducers](../02-reducers) — what happens when two nodes write the same key.
