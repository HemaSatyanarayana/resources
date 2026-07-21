# 07 — Memory: checkpointers, threads, and time travel

> **Concepts:** `InMemorySaver` · `compile(checkpointer=...)` · `thread_id` · `get_state` · `StateSnapshot` · `get_state_history` · time travel

## The idea

```python
graph = builder.compile(checkpointer=InMemorySaver())
```

After **every step**, LangGraph writes the state to the checkpointer under the
`thread_id` you pass at call time. Before every run, it loads it back. That one
line turns a stateless function into a resumable conversation.

The consequence to internalise:

```python
graph.invoke({"messages": [("user", "what is my name?")]}, thread("a"))
```

You send **one message**. The state that comes back has four. You are no longer
responsible for carrying the transcript around — and the input schema is
"what's new", not "everything so far".

### Memory is a compile-time decision

Nothing about your nodes, edges, or state schema changes. The same builder
compiles with or without memory:

| | Storage | Survives |
|---|---|---|
| no checkpointer | none | nothing |
| `InMemorySaver()` | a dict in this process | the process, not a restart |
| `SqliteSaver` / `PostgresSaver` | a database | restarts, deploys, machines |

`InMemorySaver` is for tests and labs. Everything you learn here applies
unchanged to the durable savers — swap the object, keep the code. (They ship
separately: `langgraph-checkpoint-sqlite`, `langgraph-checkpoint-postgres`.)

### `thread_id` is your partition key

Same id, same conversation. Different id, different world. It is a chat session,
a support ticket, a user id — whatever your product calls "a conversation". Two
users sharing a `thread_id` is a data leak, so derive it from something
authenticated, never from user input.

Compile with a checkpointer and *forget* the thread id and LangGraph raises.
Persistence needs a key; there is no sensible default.

### `get_state` and the snapshot

```python
snapshot = graph.get_state(thread("a"))
snapshot.values   # the state dict right now
snapshot.next     # nodes still pending — () when the run finished
snapshot.config   # includes this snapshot's checkpoint_id
```

`.next` is the field to remember. It is empty for a completed run — and *not*
empty for one that stopped in the middle, which is precisely what exercise 08
is about.

### Time travel is a consequence, not a feature

A checkpoint per step means the past is addressable. Take an old snapshot's
config, pass it to `.invoke()`, and the run continues **from there**:

```python
turn_1 = ...                                     # a snapshot from history
graph.invoke({"messages": [("user", "let's try something else")]}, turn_1.config)
# -> ["Q1", "A1", "let's try something else", "A3"]     Q2/A2 never happened
```

The branch you abandoned is still in the history; the thread's head now points
at the new one. This is what makes "undo the agent's last three steps and retry
with a hint" a normal operation rather than a rewrite.

### The bill

Every turn re-sends the whole conversation. Turn 30 costs far more than turn 1,
and eventually blows the context window. Real systems trim (drop old messages),
summarise (replace old messages with a synopsis), or store facts elsewhere and
retrieve the relevant ones. All three are ordinary nodes operating on
`state["messages"]` — you already know how to write them.

### What is written, and when

"It saves state" is too vague to reason about. Precisely: **after every
superstep** (exercise 01), LangGraph writes a checkpoint containing the values
of every channel, which node runs next, and metadata about how it got there.
Each checkpoint has a unique `checkpoint_id` and is stored under your
`thread_id`.

One turn of the one-node graph in this exercise produces **three** checkpoints,
not one. Run it and look:

```python
sizes = [len(s.values["messages"]) for s in graph.get_state_history(thread("a"))]
# [2, 1, 0]  -> after chat, before chat, and the initial write
```

Per-step rather than per-invocation is the design decision everything else rests
on. If checkpoints were written only at the end of a run:

- a crash mid-run would lose the entire run rather than one step
- you could not inspect a run that is still in progress
- you could not rewind to *within* a run — only between runs
- **exercise 08 would be impossible**, because pausing in the middle of a graph
  is precisely "checkpoint here and stop"

Human-in-the-loop is not a separate feature bolted on later. It is what you get
for free once state is durable at every step.

### `thread_id` is a partition key

Same id, same conversation. Different id, different world. It is whatever your
product calls a conversation: a chat session, a support ticket, a document being
edited, a user id for a single-session assistant.

Three things follow, in increasing order of how badly they bite:

**Compile with a checkpointer and forget the thread id, and LangGraph raises.**
Persistence needs a key; there is no sensible default and it refuses to guess.

**Two users sharing a thread id is a data leak.** One person's conversation is
loaded into another's. So derive the id from something *authenticated* —
a session the server issued, a row you own — never from a query parameter or a
client-supplied header. This is the standard IDOR mistake in a new setting, and
the fact that the key is called `thread_id` rather than `user_id` makes it easy
to forget it is a security boundary.

**Threads are not cleaned up for you.** With `InMemorySaver` they vanish with
the process. With a database they accumulate forever, holding whatever your
users said. Retention and deletion are your problem, and "delete this
conversation" means deleting rows in the checkpointer, not just hiding a UI
element.

### The `StateSnapshot`, field by field

`get_state` and `get_state_history` both hand you snapshots. Four fields matter:

| Field | What it holds | What you use it for |
|---|---|---|
| `.values` | the state dict at that point | reading the conversation |
| `.next` | node names still to run — `()` when finished | **is this run done?** |
| `.config` | thread id + this snapshot's `checkpoint_id` | addressing this exact moment |
| `.metadata` | step number, source, the writes made | debugging, auditing |

`.next` is the one to internalise. For a completed run it is empty. For a run
that stopped in the middle it names what it stopped before — which is exactly
how exercise 08 reports a pending approval:

```python
graph.get_state(config).next
# ()             -> finished
# ('approval',)  -> waiting for a human
```

An untouched thread has no state at all: `.values` is `{}`, not a schema-shaped
dict of empties. Read defensively — `snapshot.values.get("messages", [])`.

### Time travel is a consequence, not a feature

Nobody implemented "undo". It falls out of two facts you already have: every
step is checkpointed, and every checkpoint is addressable by config. So resuming
from the past is the same `.invoke()` you always call, with an older config:

```python
turn_1 = <a snapshot from history>
graph.invoke({"messages": [("user", "actually, something else")]}, turn_1.config)
```

The run continues **from that moment**, seeing only the state as it was then.

The part that surprises people is what happens to the thread afterwards. The
original branch is **not deleted** — it is still in the history — but the
thread's head now points at the **new** branch. `get_state` returns the new
line; the abandoned one is reachable only by its own checkpoint ids.

This is a branching version history, not a linear undo, and it is the mechanism
behind things that are otherwise very hard to build: "regenerate this response",
"edit my message three turns back and continue from there", "run the same
conversation against two different prompts and compare". It is also the reason
`checkpoint_id` appears in `.config` at all.

### `InMemorySaver` is not the interesting one

| Saver | Storage | Survives | Use |
|---|---|---|---|
| none | — | nothing | stateless graphs |
| `InMemorySaver` | a dict in this process | nothing | tests, labs, demos |
| `SqliteSaver` | a file | restarts | single-machine apps |
| `PostgresSaver` | a database | restarts, deploys, machines | production |

Everything you write here works unchanged against all four, because your graph
depends on `BaseCheckpointSaver`, not on a concrete class. Swapping is a
one-line change at the *edge* of your program.

Two things do change in production. **Serialisation**: everything in state must
survive a round trip through the saver, so open file handles, sockets, or
arbitrary objects will fail — keep state to data. **Concurrency**: two requests
on one thread at the same time now genuinely race, and the checkpointer is where
that plays out.

### The bill nobody warns you about

Memory is not free, and the cost is not storage — it is **tokens**. Every turn
re-sends the entire conversation (exercise 04). Turn 50 might carry 40,000
tokens of history to ask a ten-word question, and eventually the context window
simply ends the conversation.

Three standard responses, none of them provided by the checkpointer:

- **Trim** — drop old messages. Cheap, lossy, and the loss is invisible until
  someone references something you dropped. `RemoveMessage` (exercise 02) is how.
- **Summarise** — replace old turns with a synopsis. Preserves gist, costs an
  extra model call, and compounds errors as summaries get re-summarised.
- **Retrieve** — store facts elsewhere and pull in only what is relevant. Most
  work, best scaling, and the doorway to RAG.

All three are ordinary nodes operating on `state["messages"]`. Note also what
the checkpointer is *not*: it is **per-thread** state. "Remember that Hema
prefers Go, across all future conversations" is a different problem —
long-term memory, which LangGraph handles with a separate `Store` API. A
checkpointer would never surface it, because it lives under a different thread.

## The imports you need

```python
from __future__ import annotations

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import BaseMessage, SystemMessage
from langgraph.checkpoint.base import BaseCheckpointSaver
from langgraph.graph import START, MessagesState, StateGraph
```

New here:

**`BaseCheckpointSaver`** — the abstract base class every checkpointer inherits
from, used purely as a type annotation:

```python
def build_chat_graph(model, checkpointer: BaseCheckpointSaver | None = None):
```

Same injection pattern as `BaseChatModel` in exercise 04, and for the same
reason: your graph should not care *which* saver it got. `InMemorySaver` in
tests, `PostgresSaver` in production, identical code. Note it lives under
`langgraph.checkpoint`, a separate subpackage — the durable savers ship as
separate installs (`langgraph-checkpoint-sqlite`,
`langgraph-checkpoint-postgres`) and plug in here.

**`BaseMessage`** — the base class of `HumanMessage`, `AIMessage`,
`SystemMessage` and `ToolMessage`. You need it only to annotate a return type:
`list[BaseMessage]` means "a list of messages, any kind".

The concrete checkpointer is imported by the **tests**, and by your `__main__`
block, but not by your graph code:

```python
from langgraph.checkpoint.memory import InMemorySaver
```

That asymmetry is the lesson. The graph depends on the *interface*; only the
edges of your program pick an implementation.

## What to build

### 1. `SYSTEM_PROMPT`

```python
SYSTEM_PROMPT = "You are a helpful assistant with a good memory."
```

### 2. `thread(thread_id) -> dict`

Return the config dict that tells LangGraph which conversation this is:

```python
{"configurable": {"thread_id": thread_id}}
```

The nesting is not decoration — `configurable` is the runtime namespace, and
other keys live beside `thread_id` there. (`recursion_limit`, which you passed
in 05, is *not* configurable-scoped; it sits at the top level. Two different
kinds of config.)

A `thread_id` is whatever you say it is: a chat session, a ticket number, a user
id. Same id, same conversation.

### 3. `build_chat_graph(model, checkpointer=None)`

A one-node chat graph: `START -> chat -> END`. The node prepends
`SystemMessage(SYSTEM_PROMPT)` (exercise 06) and appends the model's reply. Name
the node `chat`.

The only new thing is the last line:

```python
return builder.compile(checkpointer=checkpointer)
```

`checkpointer=None` compiles a perfectly good graph with no memory — that is the
default you have been using all along. Memory is a **compile-time decision**,
not a graph-structure one: nothing about the nodes changes.

### 4. `say(graph, thread_id, text) -> str`

Send one user message on `thread_id` and return the reply's **text**.

Note what you pass: `{"messages": [("user", text)]}` — just the new message. The
checkpointer supplies the history. Do not accumulate the transcript yourself;
that is the job you are delegating.

### 5. `conversation(graph, thread_id) -> list[BaseMessage]`

Return the messages currently stored for `thread_id`.

`graph.get_state(config)` returns a `StateSnapshot`:

| Attribute | What it holds |
|---|---|
| `.values` | the state dict at this point |
| `.next` | node names still to run — empty when the run finished |
| `.config` | this snapshot's config, including its `checkpoint_id` |

A thread nobody has used yet has no state, so `.values` is `{}` — return an
empty list rather than raising.

### 6. `checkpoint_configs(graph, thread_id) -> list[dict]`

Return the config of every checkpoint on `thread_id`, newest first.

`graph.get_state_history(config)` yields a `StateSnapshot` per **step** — so one
turn of this graph produces several. Each one's `.config` carries a
`checkpoint_id`, and passing that config back to `.invoke()` resumes from
exactly that moment.

### 7. Optional: a runnable demo

```python
if __name__ == "__main__":
    from langgraph.checkpoint.memory import InMemorySaver

    from labgraph import get_chat_model, print_messages

    graph = build_chat_graph(get_chat_model(), InMemorySaver())

    print(say(graph, "demo", "My name is Hema and I like Go."))
    print(say(graph, "demo", "What is my name, and what do I like?"))

    print("\n-- stored on thread 'demo' --")
    print_messages(conversation(graph, "demo"))
    print(f"\n{len(checkpoint_configs(graph, 'demo'))} checkpoints written")
```

The test to read first is `test_a_graph_without_a_checkpointer_forgets_immediately`
— same `thread_id`, no memory. The id is a key; without a checkpointer there is
no store for it to key into.

## Run it

```bash
uv run pytest exercises/07-memory -v
uv run python exercises/07-memory/graph.py     # needs .env
```

## Think about it

- One turn of a one-node graph writes several checkpoints. Why per *step* rather
  than per *invocation*? (What would exercise 08 be unable to do otherwise?)
- Where would `thread_id` come from in a real web app — and what goes wrong if
  it comes from a query parameter?
- Swap `InMemorySaver` for `SqliteSaver` and your agent survives a restart. What
  in your state is *not* safe to persist that way?
- After time travel, both branches exist in the history. What does that let you
  build for a user who wants to compare two answers?
- The system prompt lives in code, not in state (06). You resume a year-old
  thread after changing that prompt — what happens, and is it what you want?

<details>
<summary>Reference solution</summary>

```python
from langchain_core.language_models import BaseChatModel
from langchain_core.messages import BaseMessage, SystemMessage
from langgraph.checkpoint.base import BaseCheckpointSaver
from langgraph.graph import START, MessagesState, StateGraph

SYSTEM_PROMPT = "You are a helpful assistant with a good memory."


def thread(thread_id: str) -> dict:
    return {"configurable": {"thread_id": thread_id}}


def build_chat_graph(model: BaseChatModel, checkpointer: BaseCheckpointSaver | None = None):
    def chat(state: MessagesState) -> dict:
        messages = [SystemMessage(SYSTEM_PROMPT), *state["messages"]]
        return {"messages": [model.invoke(messages)]}

    builder = StateGraph(MessagesState)
    builder.add_node("chat", chat)
    builder.add_edge(START, "chat")

    return builder.compile(checkpointer=checkpointer)


def say(graph, thread_id: str, text: str) -> str:
    result = graph.invoke({"messages": [("user", text)]}, thread(thread_id))
    return result["messages"][-1].content


def conversation(graph, thread_id: str) -> list[BaseMessage]:
    snapshot = graph.get_state(thread(thread_id))
    return snapshot.values.get("messages", [])


def checkpoint_configs(graph, thread_id: str) -> list[dict]:
    return [snap.config for snap in graph.get_state_history(thread(thread_id))]
```

`say` returns `result["messages"][-1].content` — `.invoke()` still returns the
full state, checkpointer or not. The checkpointer changes what goes *in*, not
what comes out.

</details>

---

Next: [08 — Human in the loop](../08-human-in-the-loop) — stop the graph
mid-run, ask a person, and resume.
