# 09 — Streaming: `updates`, `values`, `messages`, `custom`

> **Concepts:** `.stream()` · `stream_mode` · `AIMessageChunk` · `langgraph_node` metadata · `get_stream_writer` · multi-mode streaming

## The idea

Same graph, same answer, same cost — but the user sees it happening. Swap
`.invoke()` for `.stream()` and pick what you want to watch:

| `stream_mode` | Each chunk is | Use it for |
|---|---|---|
| `"updates"` | `{node_name: what_that_node_returned}` | "Searching…", "Calling tool…" |
| `"values"` | the whole state, after every step | live state inspectors, debugging |
| `"messages"` | `(AIMessageChunk, metadata)` per token | the typing effect |
| `"custom"` | whatever a node wrote via `get_stream_writer()` | your own progress events |

They are four views of one run, and you can subscribe to several at once:

```python
for mode, chunk in graph.stream(input, stream_mode=["updates", "custom"]):
    ...   # chunks arrive tagged with the mode that produced them
```

### `updates` vs `values`

`updates` gives you the **delta** — exactly what the node returned. `values`
gives you the **whole state**, re-emitted every step. On a 60-message
conversation, `values` copies all 60 out on every step to tell you one thing
changed. Default to `updates`; reach for `values` when you genuinely want to
render the whole thing.

### `messages` mode has two traps

```python
for chunk, metadata in graph.stream(input, stream_mode="messages"):
```

It yields **2-tuples**, and both halves matter.

**Trap 1: it is not only your model node.** Every message-producing node feeds
this stream, tool results included. Filter on
`metadata["langgraph_node"] == "model"` or your chat bubble will render `5.0` in
the middle of the sentence. The same metadata is what lets a multi-agent app
label which agent is talking.

**Trap 2: empty chunks.** A tool-calling turn streams as `AIMessageChunk`s with
no text — the payload is in `.tool_calls`, not `.content`. Skip anything falsy
or the UI flickers with empty renders.

And the invariant worth testing: **the pieces must rejoin exactly.**
`"".join(tokens) == answer`. Any stripping or space-adding while streaming shows
up as mangled text.

### `custom` is the escape hatch

```python
writer = get_stream_writer()
writer({"stage": "searching", "done": 3, "total": 10})
```

Some progress is neither state nor conversation: "3 of 10 files indexed". Do not
force it into `messages` — you would be polluting the model's context with UI
chatter, and paying tokens for it. The custom stream is a **side channel**: it
never touches state, never reaches the model, never lands in a checkpoint.

Note it fires *per node execution* — a node that runs twice in a loop emits
twice. That is usually what you want for progress.

### Streaming does not change the answer

Same graph, same final state. Streaming is a way of *observing* a run, not a
different execution mode. Which means you can develop with `.invoke()` and add
streaming at the UI layer without touching a node.

### Two different things are streaming

The four modes are easier to keep straight once you notice they come from two
distinct sources:

**Graph-level events.** LangGraph knows which node it scheduled, what that node
returned, and when the state changed. `"updates"`, `"values"` and `"custom"` are
all views of that — they exist because the *framework* is orchestrating, and
they would work identically for a graph with no LLM in it at all.

**Model-level tokens.** The provider streams the reply back a fragment at a time
over a single HTTP response. `"messages"` surfaces that inner stream, plumbed up
through the node that made the call.

That is why `"updates"` gives you one chunk for a node that took four seconds,
while `"messages"` gives you forty chunks from inside that same node. They are
answering different questions: *what is the system doing?* versus *what is the
model saying?*

A real UI usually wants both — a status line driven by `"updates"` and a text
area driven by `"messages"` — which is what multi-mode streaming is for.

### `updates` versus `values`, concretely

For the one-tool run in this exercise:

```python
# stream_mode="updates"          3 chunks — the deltas
{"model": {"messages": [AIMessage(tool_calls=[...])]}}
{"tools": {"messages": [ToolMessage("5.0")]}}
{"model": {"messages": [AIMessage("The answer is 5")]}}

# stream_mode="values"           4 chunks — the whole state, each time
{"messages": [Human]}
{"messages": [Human, AI]}
{"messages": [Human, AI, Tool]}
{"messages": [Human, AI, Tool, AI]}
```

Same run, two shapes. Note that `"values"` emits **before** any node runs — the
first chunk is your input — while `"updates"` only speaks when a node finishes.

The cost difference is not academic. `"values"` re-emits the entire state on
every step, so on a 60-message conversation you copy 60 messages out per step to
communicate that one arrived. Default to `"updates"`; reach for `"values"` when
you genuinely want to render the whole thing, or when debugging and you want to
watch state evolve.

Also note the chunk shape: `{node_name: update}`. Usually one key — but a
superstep that ran several nodes puts **all of them in one chunk**, which is why
you iterate the keys instead of taking the first. Exercise 10 makes that the
normal case.

### `messages` mode: the two-tuple and its traps

```python
for chunk, metadata in graph.stream(input, stream_mode="messages"):
```

Both halves matter, and both cause bugs.

**The chunk is an `AIMessageChunk`, not an `AIMessage`.** A chunk is a *fragment*
designed to be combined: `chunk_a + chunk_b` produces a longer chunk, and
concatenating every `.content` reproduces the final message exactly. That exact
reproduction is the invariant your UI depends on, and it is why
`stream_answer_tokens` must not strip or pad anything — a `.strip()` inside the
loop deletes the spaces *between* words, and you get `Theansweris5`.

**Trap 1: not every chunk is your model.** Every message-producing node feeds
this stream, so the `ToolMessage` from your `tools` node arrives here too —
whole, not chunked, since no model generated it. Render blindly and `5.0` lands
mid-sentence in the chat bubble. Filter on
`metadata["langgraph_node"] == "model"`.

**Trap 2: empty chunks.** A tool-calling turn streams as `AIMessageChunk`s whose
`.content` is `""` — the payload is in `.tool_calls`, not text. Emit those to a
UI and you get flickers of nothing. Skip anything falsy.

The `metadata` dict is more useful than it first appears. `langgraph_node` is
what you filter on here, but it also carries the run id and any tags you
attached — which is how a multi-agent app labels *which* agent is currently
talking, using one stream.

### `custom`: a side channel that never touches state

```python
writer = get_stream_writer()
writer({"stage": "searching", "done": 3, "total": 10})
```

Some progress is neither state nor conversation. "Indexed 3 of 10 files" is not
part of the answer, it is not something the model should ever read, and it
should not be checkpointed. It exists purely to tell a human that something is
happening.

The three places you might be tempted to put it are all wrong:

- **In `messages`** — pollutes the model's context and costs tokens forever, to
  carry UI chatter.
- **In state** — checkpointed, and now your conversation history contains
  progress bars.
- **In a log** — invisible to the user, who is the audience.

The custom stream is the fourth option: emitted, observed, and gone. It never
touches state, never reaches the model, never lands in a checkpoint.

Two mechanics worth knowing. It only works **inside a running node** —
`get_stream_writer()` needs a graph context to attach to. And it fires **per node
execution**, so a node that runs three times in a loop emits three events; the
spec's expected `[{...messages_so_far: 1}, {...messages_so_far: 3}]` is exactly
the model node running twice.

### Streaming does not change the run

Same graph, same nodes, same final state, same cost. `.stream()` is a way of
*observing* execution, not a different execution mode — which has a pleasant
practical consequence: build and test with `.invoke()`, then add streaming at
the UI layer without touching a single node. The spec asserts this directly
(`test_streaming_and_invoking_produce_the_same_final_state`).

What streaming *does* change is your operational surface, in ways worth thinking
about before production:

- **Disconnects.** The user closes the tab three seconds into a ten-second run.
  The graph does not care — it keeps running server-side. With a checkpointer
  the work is not lost and they can rejoin the thread; without one it is
  finished and discarded.
- **Errors mid-stream.** You have already sent 200 tokens when a tool fails.
  There is no unsending them, so your UI needs a way to show a failure *after*
  partial success.
- **Interrupts mid-stream.** An `interrupt` (exercise 08) ends the stream with
  the run unfinished, which your client must distinguish from normal completion
  — otherwise a paused approval looks like a truncated answer.

## The imports you need

```python
from __future__ import annotations

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import AIMessageChunk
from langgraph.config import get_stream_writer
from langgraph.graph import START, MessagesState, StateGraph
from langgraph.prebuilt import ToolNode, tools_condition

from labgraph import add, multiply, search
```

Everything is carried over from 06 except two.

**`AIMessageChunk`** — a *partial* `AIMessage`. When a model streams, you do not
get one `AIMessage`; you get a sequence of chunks, each holding a fragment of
the text. They are designed to add together (`chunk_a + chunk_b`), and joining
all their `.content` reproduces the final message exactly.

You need it for an `isinstance` check, because `stream_mode="messages"` also
emits complete `ToolMessage`s from your tools node. `isinstance(chunk,
AIMessageChunk)` is how you tell "a piece of the model's answer" from "a tool
result that arrived whole".

**`get_stream_writer`** — from `langgraph.config`, a module of helpers a node can
call to reach the machinery running it. Called **inside a node**, with no
arguments, it returns a function that pushes any object into the `"custom"`
stream:

```python
writer = get_stream_writer()
writer({"stage": "thinking"})
```

It only works inside a running node — there is no graph to talk to otherwise.
That is why you call it in the body of `call_model`, not at module level.

## What to build

### 1. `build_graph(model)`

Exercise 06's agent — `ToolNode`, `tools_condition`, nodes named `model` and
`tools` — plus **one new line** in the model node, before invoking the model:

```python
writer = get_stream_writer()
writer({"stage": "thinking", "messages_so_far": len(state["messages"])})
```

The spec pins that exact payload shape.

(No system prompt in this one — keep the model node minimal so that what you see
in each stream is only what this exercise is about.)

### 2. A small input helper

Every consumer below sends the same shape, so factor it out:

```python
def _input(text: str) -> dict:
    return {"messages": [("user", text)]}
```

### 3. `stream_node_updates(graph, text) -> list[str]`

Node names in the order they finish, using `stream_mode="updates"`.

Each chunk is `{node_name: update_it_returned}` — usually one key, but a
superstep that ran two nodes in parallel yields both, so iterate the keys rather
than taking `[0]`. For a one-tool run you should get
`["model", "tools", "model"]`.

### 4. `stream_state_sizes(graph, text) -> list[int]`

`len(state["messages"])` after each step, using `stream_mode="values"`.

Each chunk is the **entire state** — the same dict `.invoke()` would return, just
emitted repeatedly as it grows. The first chunk is the input state, before any
node has run.

### 5. `stream_answer_tokens(graph, text) -> list[str]`

The model's answer as the list of text pieces it arrived in.

`stream_mode="messages"` yields `(chunk, metadata)` **2-tuples**, not bare
chunks. Two filters are mandatory:

- `metadata["langgraph_node"] != "model"` → skip. This mode also emits the
  `ToolMessage` from the `tools` node, and rendering `"5.0"` into the chat
  bubble is the classic streaming bug.
- not an `AIMessageChunk`, or empty `.content` → skip. A tool-call turn arrives
  as chunks with no text at all.

Joined back together the pieces must reproduce the answer **exactly**, so do not
add spaces or strip anything.

### 6. `stream_progress(graph, text) -> list[dict]`

Everything the nodes wrote via `get_stream_writer()`, using
`stream_mode="custom"`. This mode yields exactly the objects you passed to the
writer, nothing else.

### 7. Optional: a runnable demo

The real thing — tokens to the terminal as they are generated:

```python
if __name__ == "__main__":
    import sys

    from labgraph import get_chat_model

    graph = build_graph(get_chat_model())
    question = "What is 128 times 4? Then explain what a checkpointer is."

    print(f"you: {question}\nmodel: ", end="", flush=True)

    for chunk, metadata in graph.stream(_input(question), stream_mode="messages"):
        if metadata.get("langgraph_node") == "model" and isinstance(chunk, AIMessageChunk):
            sys.stdout.write(str(chunk.content))
            sys.stdout.flush()
    print()
```

## Run it

```bash
uv run pytest exercises/09-streaming -v
uv run python exercises/09-streaming/graph.py     # needs .env; watch it type
```

## Think about it

- `test_updates_arrive_during_the_run_not_after` calls `next()` on the stream
  once. Why does that assertion prove streaming, when collecting the whole list
  would not?
- Tokens arrive over ten seconds and the user closes the tab at second three.
  What happens to the run, and what does the checkpointer hold? (Exercise 07.)
- The custom writer fires once per model turn. What would you emit from a node
  that loops over 500 documents — and how often?
- You want "Searching the web…" in the UI while a tool runs. Which mode, and
  where does the text come from — the node, or the UI's own mapping of node
  names?
- Streaming and the approval gate from 08: what does a stream do when the graph
  interrupts mid-run?

<details>
<summary>Reference solution</summary>

```python
from __future__ import annotations

from langchain_core.language_models import BaseChatModel
from langchain_core.messages import AIMessageChunk
from langgraph.config import get_stream_writer
from langgraph.graph import START, MessagesState, StateGraph
from langgraph.prebuilt import ToolNode, tools_condition

from labgraph import add, multiply, search

TOOLS = [add, multiply, search]


def _input(text: str) -> dict:
    return {"messages": [("user", text)]}


def build_graph(model: BaseChatModel):
    bound = model.bind_tools(TOOLS)

    def call_model(state: MessagesState) -> dict:
        writer = get_stream_writer()
        writer({"stage": "thinking", "messages_so_far": len(state["messages"])})
        return {"messages": [bound.invoke(state["messages"])]}

    builder = StateGraph(MessagesState)
    builder.add_node("model", call_model)
    builder.add_node("tools", ToolNode(TOOLS))
    builder.add_edge(START, "model")
    builder.add_conditional_edges("model", tools_condition)
    builder.add_edge("tools", "model")

    return builder.compile()


def stream_node_updates(graph, text: str) -> list[str]:
    return [
        name
        for chunk in graph.stream(_input(text), stream_mode="updates")
        for name in chunk
    ]


def stream_state_sizes(graph, text: str) -> list[int]:
    return [
        len(state["messages"])
        for state in graph.stream(_input(text), stream_mode="values")
    ]


def stream_answer_tokens(graph, text: str) -> list[str]:
    tokens = []
    for chunk, metadata in graph.stream(_input(text), stream_mode="messages"):
        if metadata.get("langgraph_node") != "model":
            continue
        if not isinstance(chunk, AIMessageChunk) or not chunk.content:
            continue
        tokens.append(chunk.content)
    return tokens


def stream_progress(graph, text: str) -> list[dict]:
    return list(graph.stream(_input(text), stream_mode="custom"))
```

`stream_node_updates` iterates `for name in chunk` rather than taking one key: a
parallel superstep puts several nodes in a single chunk. Exercise 10 makes that
the normal case.

</details>

---

Next: [10 — Parallel fan-out](../10-parallel-fanout) — one node, many workers,
`Send`.
