"""Exercise 07 — Memory: checkpointers, threads, and time travel.

Every graph so far has been amnesiac. A checkpointer saves state after every
step, keyed by a `thread_id` you supply at call time — which gives you
conversation memory, resumability, and (exercise 08) human-in-the-loop, all from
one mechanism.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/07-memory/README.md

The spec imports these names — they must exist with these exact spellings:

    SYSTEM_PROMPT       the assistant's system prompt
    thread              thread_id -> the config dict LangGraph expects
    build_chat_graph    (model, checkpointer=None) -> compiled graph
    say                 send one message on a thread, return the reply text
    conversation        the messages stored for a thread
    checkpoint_configs  every checkpoint's config, newest first

Then:

    uv run pytest exercises/07-memory -v
"""
