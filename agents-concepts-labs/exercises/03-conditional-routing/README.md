# 03 — Conditional routing: branching and looping

> **Concepts:** `add_conditional_edges` · router functions · cycles · `path_map` · `recursion_limit`

## The idea

A conditional edge is a function from state to *the name of the next node*:

```python
def route_after_review(state) -> str:
    if state["score"] >= 8:
        return END
    return "revise"

builder.add_conditional_edges("review", route_after_review)
```

After `review` finishes, LangGraph calls the router, reads the returned string,
and jumps there. That is the entire mechanism — and it is enough to express
every agent pattern you have heard of. A ReAct agent is exactly this router with
"did the model ask for a tool?" as the condition (exercise 05).

Three details that matter in practice:

**Routers are pure and cheap.** They read state and return a name. Do not call
an LLM or hit a database inside a router — do that in a node, write the result
to state, and let the router branch on it. Routers can and will be re-run.

**`path_map` decouples the return value from node names.** If your router
naturally returns `"yes"`/`"no"`, map them:

```python
builder.add_conditional_edges("review", router, {"yes": "revise", "no": END})
```

This also matters for **graph drawing**: with a bare router, LangGraph infers
possible targets from the return type annotation (hence the
`Literal["revise", "__end__"]` in this exercise). Give it a `path_map` or a
`Literal` and your diagrams stay accurate.

**A router that returns a node you already visited creates a cycle.** That is
allowed and it is the point — but a cycle needs a **budget**, or it never ends.

### Two different stopping conditions

This exercise deliberately has both, and confusing them is a common bug:

| | What it is | Who owns it |
|---|---|---|
| `score >= 8` | The *happy* exit — the work is done. | Your logic |
| `revisions >= MAX_REVISIONS` | The *budget* exit — out of attempts. | Your logic |
| `recursion_limit` | A hard ceiling on supersteps; raises `GraphRecursionError`. | The framework |

`recursion_limit` (default 25) is a backstop against a bug, not a control-flow
tool. If your agent routinely hits it, the fix is a budget in state, not a
bigger limit. Notice the shape of the "ai" run in the tests: it exits *unhappy*,
having spent its budget without ever passing review. Real agents need that path
to exist, and need to report it — an agent that silently returns a bad answer
after burning its budget is worse than one that says it failed.

```
START ──▶ draft ──▶ review ──┬──(score>=8 or out of budget)──▶ END
                     ▲       │
                     └──revise◀──(otherwise)
```

### How LangGraph knows where a router can go

A conditional edge poses a problem the framework cannot solve by reading your
code: your router is an arbitrary Python function, so its possible return values
are undecidable until it runs. But LangGraph wants to know the *possible*
targets ahead of time — to validate the graph at compile time, and to draw it.

There are three ways to tell it, in increasing order of explicitness:

**1. Infer from the return annotation.** This is why the exercise annotates
`-> Literal["revise", "__end__"]`. LangGraph reads that annotation and learns
the branch has two possible destinations.

**2. Pass a `path_map`.** A dict translating return values into node names:

```python
builder.add_conditional_edges("review", router, {"yes": "revise", "no": END})
```

Now the router can return domain vocabulary — `"yes"`, `"escalate"`,
`"needs_human"` — instead of node names. This decouples *the decision* from
*the topology*: rename a node and you change one dict entry, not the logic.

**3. Pass a list of possible targets.** `["research", "write"]` — used in
exercise 10, where the router returns `Send` objects that no annotation could
describe.

Give it none of the three and your graph still runs perfectly. What breaks is
the *drawing*: `get_graph()` shows a node with no visible outgoing branch. Your
diagram silently stops matching your agent, which is a documentation bug of the
worst kind — the one that looks fine.

### Supersteps, cycles, and what `recursion_limit` actually counts

Exercise 01 introduced the **superstep**: one round of "run the scheduled
node(s), merge their updates, work out what runs next". A cycle simply means a
later superstep schedules a node an earlier one already ran.

Trace the `"prompts"` topic (7 characters → score 7):

| Superstep | Node | After merging | Router says |
|---|---|---|---|
| 1 | `draft` | score 7, revisions 0 | — |
| 2 | `review` | history +review | 7 < 8, budget fine → `revise` |
| 3 | `revise` | score 8, revisions 1 | — |
| 4 | `review` | history +review | 8 ≥ 8 → `END` |

Four supersteps, two of them running the *same* node. Nothing in the graph is
"looping" in the Python sense; there is no `while`. Each superstep is an
independent scheduling decision made from state alone — which is exactly why the
whole run can be checkpointed and resumed between any two of them (exercise 07).

**`recursion_limit` counts supersteps, not visits to a particular node.** The
default is 25. When exceeded, LangGraph raises `GraphRecursionError`. Two things
follow:

- A graph with many nodes per pass burns the budget faster than its loop count
  suggests. Four supersteps per "iteration" means the default allows about six
  iterations, not 25.
- It is a **circuit breaker**, not a control-flow mechanism. Raising the limit
  because you hit it is almost always the wrong fix — it converts a fast failure
  into a slow, expensive one.

### Designing termination properly

This is the part that separates a demo from something you would run for real,
and it is why the exercise has three stopping conditions rather than one.

**The happy exit** (`score >= 8`) means the work succeeded.

**The budget exit** (`revisions >= MAX_REVISIONS`) means you ran out of
attempts. Note where the counter lives: **in state**, not in a module-level
variable. Three reasons, each fatal on its own:

1. A module variable is shared by every concurrent run in the process — two
   users' agents would consume each other's budget.
2. It is not checkpointed, so a resumed run (exercise 07) forgets how much it
   already spent.
3. Rewinding to an earlier checkpoint would not rewind the counter, so time
   travel produces nonsense.

Anything a decision depends on belongs in state. That is not a LangGraph
convention; it is what makes the run reconstructible.

**The framework backstop** (`recursion_limit`) means *you have a bug*.

Now the part people skip. Look at what the `"ai"` run does: score 2, +1 per
revision, budget 3 — it can never reach 8, so it exits **unhappy**, having spent
its whole budget. Your graph must have that path, and the caller must be able to
tell the two exits apart:

```python
if result["score"] >= 8:  ...   # succeeded
else:                     ...   # gave up
```

An agent that burns its budget and returns its best bad answer *as if it had
succeeded* is worse than one that crashes, because the failure is now silent and
downstream. Deciding how a run reports "I could not do this" is a design
decision you own — the framework will not make it for you.

### Why routers must stay pure and cheap

A router is a function of state, and it will be called on **every pass** through
its source node. Two rules follow:

**No side effects.** Do not call an LLM, hit a database, or write a file inside a
router. If the decision needs data, fetch it in a *node*, write it to state, and
let the router branch on the value. The exercise models this exactly: `review`
is an almost-empty node whose only job is to exist so an edge can hang off it,
and `route_after_review` reads `score`. That split looks like ceremony until you
need to know *why* a run branched — at which point the reason is sitting in
state, checkpointed, streamable, and testable in isolation.

**No expense.** Routers run more often than nodes do in a loop. Keep them to
comparisons.

And a subtle one: a router returning a node name is *not* a function call. It
returns a **string**. `return END` works because `END` is the string
`"__end__"`; `return "END"` would name a node that does not exist and fail at
run time, not at compile time.

## The imports you need

```python
from __future__ import annotations

import operator
from typing import Annotated, Literal, TypedDict

from langgraph.graph import END, START, StateGraph
```

Everything here is carried over from 01 and 02 except one:

**`Literal`** — a type meaning "one of these exact values":

```python
def route_after_review(state) -> Literal["revise", "__end__"]:
```

This is not decoration. LangGraph **reads the return annotation** to work out
which nodes a conditional edge can reach, because it cannot see inside your
function body. That inferred edge list is what `get_graph()` draws and what the
spec asserts on. Omit the `Literal` and your graph still runs correctly, but its
rendered diagram loses the branch — a quiet documentation bug.

Why `"__end__"` and not `END` inside the `Literal`? `END` *is* the string
`"__end__"`, and a `Literal` needs literal values, not references to constants.
So you annotate with the raw string and `return END` in the body.

**`operator` / `Annotated`** are back for `history`, an accumulating log —
exactly exercise 02's `operator.add` pattern. Watch this key: it is how the
tests prove which path the graph actually took, and the cheapest debugging tool
you have for a loop.

## What to build

### 1. `MAX_REVISIONS`

A module-level constant, `3`.

### 2. `DraftState`

| Key | Type | Meaning |
|---|---|---|
| `topic` | `str` | what to write about |
| `draft` | `str` | the current draft |
| `score` | `int` | quality of the current draft, 0–10 |
| `revisions` | `int` | how many times we have revised |
| `history` | `Annotated[list[str], operator.add]` | which nodes ran, in order |

### 3. Three nodes

- **`write_draft(state) -> dict`** — the first draft (a stand-in for an LLM
  call). Set:
  - `draft` → `f"draft of {topic}"`
  - `score` → `len(topic) % 11` (deterministic pseudo-quality, 0–10)
  - `revisions` → `0`
  - `history` → `["draft"]`

- **`revise_draft(state) -> dict`** — each revision bumps the score by 1,
  capped at 10. Set:
  - `draft` → `f"revision {n} of {topic}"`, where `n` is the **new** revision count
  - `score` → `min(10, current score + 1)`
  - `revisions` → current + 1
  - `history` → `["revise"]`

- **`review(state) -> dict`** — a pass-through that only records that review
  happened: return `{"history": ["review"]}`. Nodes are allowed to be this
  boring; the interesting part is the edge hanging off them.

The scores are deterministic so the tests can pin exact loop counts:

- `"langgraph"` → 9 → passes immediately, 0 revisions
- `"prompts"` → 7 → one revision → 8 → passes
- `"ai"` → 2 → can never reach 8 → stops at the budget

### 4. `route_after_review(state) -> Literal["revise", "__end__"]`

**The decision.** Return `"revise"` to loop, or `END` to stop. Stop when
**either**:

- the draft is good enough — `score >= 8`, **or**
- you have already revised `MAX_REVISIONS` times (budget exhausted).

Return a node *name*, not a node. And return the `END` sentinel, not the string
`"END"` — that would be a node name that does not exist.

### 5. `build_graph()`

```
START -> draft -> review
review -> (conditional) -> revise | END
revise -> review
```

Node names: `draft`, `review`, `revise`. The conditional edge is
`builder.add_conditional_edges("review", route_after_review)`.

Note `revise` gets a *plain* edge back to `review`, not a conditional one. Only
one node in a cycle needs to make a decision; the rest just flow.

### 6. Optional: a runnable demo

```python
if __name__ == "__main__":
    from labgraph import print_state

    graph = build_graph()
    for topic in ("ai", "langgraph state machines"):
        result = graph.invoke(
            {"topic": topic, "draft": "", "score": 0, "revisions": 0, "history": []}
        )
        print_state(result, title=f"topic={topic!r}")
```

## Run it

```bash
uv run pytest exercises/03-conditional-routing -v
uv run python exercises/03-conditional-routing/graph.py
```

## Think about it

- Why must the budget counter live in *state* rather than a module-level variable?
- What breaks if `revise` forgets to increment `revisions`?
- The router returns `END`, a sentinel — what would returning `"END"` (a string) do?
- How would you tell the caller *why* a run stopped — passed review vs. gave up?
- Could `route_after_review` return a *list* of node names? What would that mean?
  (Exercise 10.)

<details>
<summary>Reference solution</summary>

```python
from __future__ import annotations

import operator
from typing import Annotated, Literal, TypedDict

from langgraph.graph import END, START, StateGraph

MAX_REVISIONS = 3


class DraftState(TypedDict):
    topic: str
    draft: str
    score: int
    revisions: int
    history: Annotated[list[str], operator.add]


def write_draft(state: DraftState) -> dict:
    return {
        "draft": f"draft of {state['topic']}",
        "score": len(state["topic"]) % 11,
        "revisions": 0,
        "history": ["draft"],
    }


def revise_draft(state: DraftState) -> dict:
    revisions = state["revisions"] + 1
    return {
        "draft": f"revision {revisions} of {state['topic']}",
        "score": min(10, state["score"] + 1),
        "revisions": revisions,
        "history": ["revise"],
    }


def review(state: DraftState) -> dict:
    return {"history": ["review"]}


def route_after_review(state: DraftState) -> Literal["revise", "__end__"]:
    if state["score"] >= 8:
        return END
    if state["revisions"] >= MAX_REVISIONS:
        return END
    return "revise"


def build_graph():
    builder = StateGraph(DraftState)

    builder.add_node("draft", write_draft)
    builder.add_node("review", review)
    builder.add_node("revise", revise_draft)

    builder.add_edge(START, "draft")
    builder.add_edge("draft", "review")
    builder.add_conditional_edges("review", route_after_review)
    builder.add_edge("revise", "review")

    return builder.compile()
```

Note `revise` has a plain edge back to `review`, not a conditional one. Only one
node in the cycle needs to make a decision; the rest just flow.

</details>

---

Next: [04 — An LLM in a node](../04-llm-node) — the first non-deterministic step.
