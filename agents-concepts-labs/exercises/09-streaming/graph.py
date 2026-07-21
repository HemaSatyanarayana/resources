"""Exercise 09 — Streaming: updates, values, messages, custom.

Your agent takes ten seconds and shows nothing until it is done. Same work, same
answer — but a spinner feels broken and text appearing as it is generated feels
fast. This is the cheapest UX win in the course.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/09-streaming/README.md

The spec imports these names — they must exist with these exact spellings:

    build_graph            model -> compiled graph that emits progress events
    stream_node_updates    node names, in finishing order   (stream_mode="updates")
    stream_state_sizes     message count after each step    (stream_mode="values")
    stream_answer_tokens   the answer, in the pieces it arrived in ("messages")
    stream_progress        your own progress events         (stream_mode="custom")

Then:

    uv run pytest exercises/09-streaming -v
"""
