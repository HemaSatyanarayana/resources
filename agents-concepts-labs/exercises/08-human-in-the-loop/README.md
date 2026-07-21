# 08 — Human in the loop: `interrupt()` and `Command(resume=...)`

> **Concepts:** `interrupt()` · `Command(resume=...)` · `__interrupt__` · `snapshot.next` · approval gates · nodes that run twice

## The idea

An autonomous agent is a liability the first time a tool does something
irreversible. The fix is not a smarter prompt — it is a **gate**: stop before
the dangerous step, ask a person, continue with their answer.

```python
answer = interrupt({"question": "Approve?", "tool_calls": [...]})
```

Three things happen at that line:

1. LangGraph **checkpoints** the run.
2. `.invoke()` **returns**, with the payload under `result["__interrupt__"]`.
3. The run is **frozen in the checkpointer**, waiting.

Later — a second later, a week later, from a different process after a deploy —
someone answers:

```python
graph.invoke(Command(resume=True), config)
```

and the graph picks up from where it stopped. Not from the top: the earlier
model call is not repeated.

### The one thing that trips everyone up

**`interrupt()` does not return a value on the first pass. It raises.**

The node did not "block". It was abandoned and rolled back to its last
checkpoint. On resume, **the whole node runs again from the top**, and this time
`interrupt()` returns what you passed to `Command(resume=...)`.

So a node containing an interrupt executes twice. Anything above the
`interrupt()` line happens twice:

```python
def review(state):
    charge_card()                    # ❌ charged twice
    answer = interrupt({...})
    ...
```

Keep side effects below the interrupt, or in a different node. This is the same
discipline that makes retries safe, and it is why interrupts and checkpoints are
the same feature.

### The pause is state, not a callback

There is no open connection, no thread parked in memory, no `await` holding a
process alive. The paused run is **rows in the checkpointer**, addressed by
`thread_id`. Which means:

- `graph.get_state(config).next` is `("approval",)` instead of `()` — exercise
  07's snapshot, now telling you the run is unfinished and where it stopped.
- Two threads can sit pending at once, independently.
- Approving one approves exactly one.
- A restart loses nothing (with a durable saver).

That is what makes this usable as a real approval queue: pending reviews are
just rows, and your UI is a list of them.

### Denial has to be a tool result

A denial is not "end the run". It is an **observation the model has to see**:

```python
ToolMessage(f"Denied by a human: {name} was not executed.", tool_call_id=call["id"])
```

Two independent reasons this must exist. **Protocol:** providers reject a
conversation containing a tool call with no matching result — a dangling call is
an API error on the very next request. **Behaviour:** a model that does not
learn its request was refused will simply ask again, and now you have an
approval loop that a tired human clicks through.

Answer *every* call on that message, not just the sensitive ones.

### Gate the minimum

`SENSITIVE_TOOLS` is a set, not a flag on the whole agent. If every tool needs
approval, reviewers stop reading and click approve — a rubber stamp is worse
than no gate, because it manufactures the appearance of oversight. Gate what is
irreversible, expensive, or public.

### The interrupt lifecycle, step by step

This is the mechanism people most often get wrong, so here it is in full.

**First pass.**

1. The graph reaches `approval` and runs `review_tool_calls(state)`.
2. Your code calls `interrupt({...})`.
3. `interrupt` **raises** a `GraphInterrupt` exception. It does not return.
4. LangGraph catches it. The node's work is **discarded** — nothing it computed
   is kept, because it never returned an update.
5. The graph checkpoints, with `.next == ("approval",)`: "the next thing to run
   is `approval`, and it has not run yet."
6. `.invoke()` **returns** to your caller, with the payload under
   `result["__interrupt__"]`.

At this point there is no thread parked in memory, no open connection, no
coroutine awaiting anything. The run is **rows in the checkpointer**. Your
process can restart, deploy a new version, and scale to zero.

**Second pass — someone answers.**

7. `graph.invoke(Command(resume=True), config)` — note there is no state dict.
   `Command` is an instruction, not input.
8. LangGraph loads the checkpoint and sees `approval` pending.
9. **`review_tool_calls` runs again, from the very first line.**
10. This time, when execution reaches `interrupt(...)`, it **returns `True`**
    instead of raising.
11. The node returns normally, the update merges, and the graph carries on.

### The consequence: your node runs twice

Everything above the `interrupt()` call executes on both passes.

```python
def review(state):
    charge_card()                # ❌ charged twice
    answer = interrupt({...})
    ...
```

This is not a quirk to work around; it is the same rule that governs retries and
resumed runs (exercises 01 and 07), just made unavoidable. Nodes must be safe to
re-execute. Practical discipline:

- **Read, compute, and format above the interrupt.** Those are idempotent.
- **Put side effects below it**, or in a different node — which is exactly why
  the refund executes in `tools`, after the gate, and not inside `approval`.
- **Never mint an id, timestamp, or random value above an interrupt** and expect
  it to survive: the second pass generates a different one.

The spec pins this: `test_resuming_does_not_replay_the_model_call` asserts the
model was called once before the pause and once after. The *node* re-runs; the
graph does not restart. Everything before `approval` stays done.

### The pause is state, not a callback

Compare the two mental models, because the wrong one leads to a design that
cannot work:

| | Callback model (wrong) | Checkpoint model (actual) |
|---|---|---|
| What is waiting | a thread/coroutine in memory | rows in a database |
| Survives a restart | no | yes |
| Approve from another process | no | yes |
| 10,000 pending approvals cost | 10,000 live handles | 10,000 rows |
| How you find pending work | a registry you built | query the checkpointer |

Because a pause is just state, everything from exercise 07 keeps working on it.
`get_state(config).next` tells you *what* it is waiting on. Threads stay
independent, so many approvals can be outstanding at once and answering one
answers exactly one. And an approval queue UI is a list of threads whose `.next`
is non-empty — no separate bookkeeping.

### The payload is the contract with the human

```python
interrupt({"question": "Approve these tool calls?",
           "tool_calls": [{"name": c["name"], "args": c["args"]} for c in sensitive]})
```

Whatever you pass comes back to the caller verbatim, so it must contain
**everything a reviewer needs to decide** — the tool, the arguments, and enough
context to judge them. A payload of `"Approve?"` produces rubber-stamping,
because the reviewer has nothing to review.

The return value is equally open. `Command(resume=X)` delivers any `X` your node
can interpret, so approval is not limited to yes/no:

```python
Command(resume=True)                    # approve
Command(resume=False)                   # deny
Command(resume={"amount": 10})          # approve, but change the arguments
```

The last one is the interesting one, and the exercise deliberately leaves it as
a *Think about it*: a reviewer who can only accept or reject will reject things
that were nearly right.

### Denial has to be a tool result — for two independent reasons

```python
ToolMessage(f"Denied by a human: {name} was not executed.", tool_call_id=call["id"])
```

**Protocol.** Every `tool_call` requires a matching `ToolMessage` (exercise 05).
Leave one dangling and the *next* model request is an API error. This is why you
answer **every** call on that message, not just the sensitive ones — the
non-sensitive calls in the same message are equally dangling.

**Behaviour.** A model that does not learn its request was refused simply asks
again. You get an approval loop, a fatigued human, and eventually an approval
that should not have been given. Telling the model what happened lets it
apologise, ask a clarifying question, or find another route — which is why
denial routes back to `model` rather than to `END`.

### Designing the gate

**Gate the minimum.** `SENSITIVE_TOOLS` is a set of names, not a switch on the
whole agent. Gate everything and reviewers stop reading and start clicking, which
is worse than no gate: you now have the *appearance* of oversight, and an audit
trail claiming a human approved something nobody looked at. Gate what is
irreversible, expensive, or publicly visible.

**Make the gate unskippable by construction.** The routing is
`model → approval → tools`, with no edge from `model` to `tools`. There is no
path around the check — not because a condition forbids it, but because the
topology has no such edge. That is far stronger than an `if` inside `tools` that
a refactor can drop, and it is visible in the rendered diagram, which makes it
reviewable.

**Know what this does not give you.** Nothing here records *who* approved, or
when, or what they saw. Real approval workflows need that, and it is a design
question with a real trade-off: identity in state means it is checkpointed and
auditable, but also that it is copied into every fork of the thread. Nothing
here handles an approval that never arrives, either — pending runs accumulate
silently, so someone has to sweep them.

## The imports you need

```python
from __future__ import annotations

from typing import Literal

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import ToolMessage
from langchain_core.tools import tool
from langgraph.checkpoint.base import BaseCheckpointSaver
from langgraph.graph import END, START, MessagesState, StateGraph
from langgraph.prebuilt import ToolNode
from langgraph.types import Command, interrupt

from labgraph import search
```

Two new names, and one new module path.

**`@tool`** — the decorator that turns a plain function into a LangChain tool.
You have been *using* tools since 05; now you write one:

```python
@tool
def refund(customer: str, amount: float) -> str:
    """Refund money to a customer. Spends real money."""
    ...
```

The decorator reads the **signature** to build the argument schema the model
must fill in, and the **docstring** to tell the model what the tool is for. Both
are sent to the provider. A tool with a vague docstring is a tool the model
misuses — treat it as prompt engineering, because it is.

**`interrupt`** — from `langgraph.types`, the module holding LangGraph's runtime
control-flow objects (as opposed to `langgraph.graph`, which holds the building
blocks). Call it inside a node to stop the run and surface a payload.

**`Command`** — from the same module. `Command(resume=value)` is what you pass
to `.invoke()` *instead of* a state dict, to answer a pending interrupt. Your
graph code does not construct one; the **caller** does. You need the import only
for the optional demo — and it is worth importing there deliberately, so you see
that resuming is something done to the graph from outside.

`tools_condition` is **not** imported this time. Your routing is no longer "tools
or end" — it is "approval or end", then "tools or model". The prebuilt router
does not fit, which is exactly the moment exercise 06 told you to drop back to
writing your own.

## What to build

### 1. The dangerous tool and the ledger

```python
LEDGER: list[dict] = []


@tool
def refund(customer: str, amount: float) -> str:
    """Refund money to a customer. Spends real money."""
    LEDGER.append({"customer": customer, "amount": amount})
    return f"Refunded ${amount:.2f} to {customer}."


TOOLS = [search, refund]
SENSITIVE_TOOLS = {"refund"}
```

`LEDGER` stands in for a payments API — the tests assert on it to prove a denied
refund really never happened. `SENSITIVE_TOOLS` is a **set of names**, not a flag
on the whole agent: everything else runs unattended.

### 2. `ApprovalState`

`MessagesState` plus `decision: str`. That key is how the approval node talks to
the router — the same "node writes a value, edge reads it" split as exercise 04.

### 3. `make_model_node(model)`

The usual model node, bound to `TOOLS`. Unchanged from 05.

### 4. `review_tool_calls(state) -> dict`

Pause for a human when the model asks for something sensitive.

Read `.tool_calls` off the last message and pick out the ones whose name is in
`SENSITIVE_TOOLS`.

**Nothing sensitive** → return `{"decision": "approved"}` and the run continues
without ever stopping. (The cheapest gate is the one that does not fire.)

**Something sensitive** → call `interrupt(payload)`, where payload is a dict a
reviewer could act on:

```python
{"question": "Approve these tool calls?",
 "tool_calls": [{"name": c["name"], "args": c["args"]} for c in sensitive]}
```

Treat `True` or `"approve"` as approval; anything else is a denial.

- On approval: `{"decision": "approved"}`.
- On denial: `{"decision": "denied", "messages": [...]}` where the messages are
  one `ToolMessage` per call on that AI message — **every** call, not just the
  sensitive ones — with content saying it was **denied** and naming the tool.

### 5. `route_after_model(state) -> Literal["approval", "__end__"]`

Tool calls go to `approval` (**not** straight to `tools`). Otherwise `END`.

### 6. `route_after_approval(state) -> Literal["tools", "model"]`

`"approved"` → `"tools"`. Anything else → back to `"model"`.

Denial does not end the run: the model gets another turn, sees the denial, and
can apologise, ask a question, or try a different approach.

### 7. `build_graph(model, checkpointer=None)`

```
START    -> model
model    -> (route_after_model)    -> approval | END
approval -> (route_after_approval) -> tools | model
tools    -> model
```

Node names: `model`, `approval`, `tools`. Use `ToolNode(TOOLS)`. Compile with
the checkpointer — an interrupt without one has nowhere to store the paused run.
**The pause is a checkpoint.**

Note the routing: `model -> approval -> tools`, never `model -> tools`. The gate
is unskippable *by construction* — no node has an edge that goes around it. That
is worth more than a check inside `tools` that someone can refactor away.

### 8. Optional: a runnable demo

```python
if __name__ == "__main__":
    from langgraph.checkpoint.memory import InMemorySaver

    from labgraph import get_chat_model, print_messages

    config = {"configurable": {"thread_id": "demo"}}
    graph = build_graph(get_chat_model(), InMemorySaver())

    result = graph.invoke(
        {"messages": [("user", "Refund $20 to customer ana, she was double charged.")]},
        config,
    )

    if "__interrupt__" in result:
        print("PAUSED — needs approval:")
        print(f"  {result['__interrupt__'][0].value}")
        answer = input("approve? [y/N] ").strip().lower() == "y"
        result = graph.invoke(Command(resume=answer), config)

    print_messages(result["messages"])
    print(f"\nledger: {LEDGER}")
```

## Run it

```bash
uv run pytest exercises/08-human-in-the-loop -v
uv run python exercises/08-human-in-the-loop/graph.py     # needs .env; prompts you
```

## Think about it

- `test_resuming_does_not_replay_the_model_call` asserts `call_count == 1` after
  the pause. What would a *re-run from the top* cost you in money and in
  duplicated side effects?
- The reviewer can only say yes or no. How would you let them *edit* the
  arguments before approving — refund $10 instead of $20? (What does the node
  do with a `Command(resume={"amount": 10})`?)
- Nothing here checks *who* approved. Where would identity go, and would you put
  it in state?
- A run pauses and nobody ever answers. How do you find it, and what should
  happen to it?
- `interrupt()` works anywhere, not just before tools. What else is worth
  pausing on — a plan before it executes? a draft before it sends?

<details>
<summary>Reference solution</summary>

```python
from __future__ import annotations

from typing import Literal

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import ToolMessage
from langchain_core.tools import tool
from langgraph.checkpoint.base import BaseCheckpointSaver
from langgraph.graph import END, START, MessagesState, StateGraph
from langgraph.prebuilt import ToolNode
from langgraph.types import interrupt

from labgraph import search

LEDGER: list[dict] = []


@tool
def refund(customer: str, amount: float) -> str:
    """Refund money to a customer. Spends real money."""
    LEDGER.append({"customer": customer, "amount": amount})
    return f"Refunded ${amount:.2f} to {customer}."


TOOLS = [search, refund]
SENSITIVE_TOOLS = {"refund"}


class ApprovalState(MessagesState):
    decision: str


def make_model_node(model: BaseChatModel):
    bound = model.bind_tools(TOOLS)

    def call_model(state: ApprovalState) -> dict:
        return {"messages": [bound.invoke(state["messages"])]}

    return call_model


def review_tool_calls(state: ApprovalState) -> dict:
    calls = state["messages"][-1].tool_calls
    sensitive = [c for c in calls if c["name"] in SENSITIVE_TOOLS]

    if not sensitive:
        return {"decision": "approved"}

    answer = interrupt(
        {
            "question": "Approve these tool calls?",
            "tool_calls": [{"name": c["name"], "args": c["args"]} for c in sensitive],
        }
    )

    if answer is True or answer == "approve":
        return {"decision": "approved"}

    return {
        "decision": "denied",
        "messages": [
            ToolMessage(
                content=f"Denied by a human: {c['name']} was not executed.",
                tool_call_id=c["id"],
                name=c["name"],
            )
            for c in calls
        ],
    }


def route_after_model(state: ApprovalState) -> Literal["approval", "__end__"]:
    if getattr(state["messages"][-1], "tool_calls", None):
        return "approval"
    return END


def route_after_approval(state: ApprovalState) -> Literal["tools", "model"]:
    if state["decision"] == "approved":
        return "tools"
    return "model"


def build_graph(model: BaseChatModel, checkpointer: BaseCheckpointSaver | None = None):
    builder = StateGraph(ApprovalState)

    builder.add_node("model", make_model_node(model))
    builder.add_node("approval", review_tool_calls)
    builder.add_node("tools", ToolNode(TOOLS))

    builder.add_edge(START, "model")
    builder.add_conditional_edges("model", route_after_model)
    builder.add_conditional_edges("approval", route_after_approval)
    builder.add_edge("tools", "model")

    return builder.compile(checkpointer=checkpointer)
```

`review_tool_calls` returns before interrupting when nothing is sensitive. The
cheapest gate is the one that does not fire.

</details>

---

Next: [09 — Streaming](../09-streaming) — the agent works; now stop making the
user stare at a spinner while it does.
