# agents-concepts-labs

Hands-on LangGraph labs: **build agents from the graph up.**

Ten exercises. Each one gives you a **README that teaches the concepts and
explains every import**, a **spec** (a test suite), and an **empty `graph.py`**.
You write the file from scratch — imports first — until the tests go green.

Nothing is pre-filled. There are no function signatures waiting for a body, and
no `TODO` markers to fill in. Deciding what to import, what the state schema
looks like, and how the pieces fit together *is* the exercise; typing a body
into someone else's skeleton teaches far less than it appears to.

No exercise asks you to trust a black box. You hand-write the ReAct loop before
you are allowed to use `create_agent`, so when the prebuilt one misbehaves you
know exactly which part is lying to you.

## Setup

```bash
uv sync
uv run pytest exercises/01-state-graph -v
```

That is the whole setup. **The tests need no API key and no network** — they run
against `labgraph.ScriptedChatModel`, a fake that replays a fixed list of
replies. The suite is deterministic and offline by design.

To also run a graph against a real model, copy `.env.example` to `.env` and fill
in three variables. Any OpenAI-compatible provider works (OpenRouter, Gemini's
compat endpoint, a local Ollama) — no provider is hardcoded anywhere in the
course. Then:

```bash
uv run python exercises/05-tool-calling/graph.py
```

The `@pytest.mark.skipif(not llm_configured())` test at the bottom of each spec
starts running too, once `.env` exists.

## The exercises

Work them in order. Each builds on the last, and most of them are one new idea.

| | Exercise | What you learn |
|---|---|---|
| 01 | [state-graph](exercises/01-state-graph) | `StateGraph`, nodes, edges, `TypedDict` state |
| 02 | [reducers](exercises/02-reducers) | how updates merge — `operator.add`, `add_messages`, your own |
| 03 | [conditional-routing](exercises/03-conditional-routing) | branching, cycles, budgets, `recursion_limit` |
| 04 | [llm-node](exercises/04-llm-node) | an LLM inside a node; injecting the model; parsing output |
| 05 | [tool-calling](exercises/05-tool-calling) | `bind_tools`, `ToolMessage`, **the ReAct loop, by hand** |
| 06 | [prebuilt-agent](exercises/06-prebuilt-agent) | `ToolNode`, `tools_condition`, `create_agent`, system prompts |
| 07 | [memory](exercises/07-memory) | checkpointers, `thread_id`, `get_state`, time travel |
| 08 | [human-in-the-loop](exercises/08-human-in-the-loop) | `interrupt()`, `Command(resume=...)`, approval gates |
| 09 | [streaming](exercises/09-streaming) | `stream_mode`: updates, values, messages, custom |
| 10 | [parallel-fanout](exercises/10-parallel-fanout) | `Send`, dynamic parallelism, map-reduce |

**01–04** are the graph itself: state in, state out, and how a graph decides
where to go. Nothing here is about AI.

**05–06** are the agent: a model that asks for tools, and a loop that runs them.

**07–10** are what turns a working loop into a working product: memory,
oversight, responsiveness, and throughput.

## How to work an exercise

1. Read the exercise's `README.md` — it is the lesson. Four sections matter:
   **The idea** (the concepts, worked from first principles), **The imports you
   need** (every import explained: what it is, where it comes from, why this
   exercise wants it), **What to build** (the exact names, signatures and
   semantics the spec expects), and **Think about it**.
2. Open `graph.py`. It holds a short brief and nothing else — the list of names
   the spec will import. You write everything below it.
3. Read the test file. **It is the spec**, and it is more precise than prose.
   When the README and the test seem to disagree, the test wins.
4. Write the file, from the imports down, until green:

   ```bash
   uv run pytest exercises/05-tool-calling -v
   ```

   While `graph.py` is still empty you get a **collection error**, not failures
   — there is nothing to import yet. Write the imports and the first name, and
   it turns into ordinary red tests.

5. Run it for real against a model, if you have `.env` set up:

   ```bash
   uv run python exercises/05-tool-calling/graph.py
   ```

6. Only then open the **Reference solution** at the bottom of the README, and
   diff it against what you wrote. Different is often fine; understand *why* it
   is different.

Each README ends with a **Think about it** section. Those questions have no
tests. They are the difference between passing the exercise and being able to
design one of these yourself — several of them are the actual design decisions
you will face in production.

`uv run pytest` with no arguments runs all ten at once, if you want to see
overall progress. Expect a wall of red until you start.

## What is already written for you

`labgraph/` is scaffolding, not an exercise. Read it — it is short.

| | |
|---|---|
| `labgraph.get_chat_model()` | builds a model from `.env`; provider-agnostic |
| `labgraph.ScriptedChatModel` | the deterministic fake: replays scripted replies, records prompts, supports `bind_tools` and streaming |
| `labgraph.tool_call` / `tool_calls` | build an assistant reply that requests tools |
| `labgraph.search / add / multiply / word_count` | toy tools for the agent labs |
| `labgraph.print_messages` / `print_state` | readable output when you run a graph by hand |
| `labgraph.TODO` | optional: raise it from a stub you have not written yet, so the module still imports while you work on one piece at a time |

`ScriptedChatModel` is worth understanding early. Because the model is
*injected* into every `build_graph(model)`, the same graph runs against a real
provider in production and a scripted fake in tests. That is why this suite is
fast, free, and deterministic — and it is a pattern to copy into your own work.

## Where to go after 10

The natural next steps, roughly in order of how often you will want them:

- **Subgraphs** — `add_node("researcher", other_compiled_graph)`. A node can be
  a whole graph, which is how exercise 10's workers become full agents.
- **Multi-agent** — supervisor and swarm patterns, built out of 03's routing and
  10's `Send`.
- **Long-term memory** — the `Store` API: facts that outlive a `thread_id`,
  unlike the checkpointer's per-conversation state.
- **Context engineering** — trimming and summarising `messages` before the
  window (and the bill) gets away from you.
- **Deployment** — LangGraph Server, and swapping `InMemorySaver` for
  `PostgresSaver`.
- **Evaluation and tracing** — LangSmith, once "did my change make it better?"
  stops being answerable by reading the output.
