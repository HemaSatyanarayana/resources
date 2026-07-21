"""Exercise 06 — The prebuilt agent: ToolNode, tools_condition, create_agent.

Exercise 05 was ~40 lines of loop. This exercise deletes almost all of it: first
by swapping your two hand-written pieces for the prebuilt ones, then by
replacing the whole graph with a single `create_agent` call.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/06-prebuilt-agent/README.md

The spec imports these names — they must exist with these exact spellings:

    TOOLS                 the list of tools the agent may use
    SYSTEM_PROMPT         the agent's system prompt
    build_graph           hand-wired, on ToolNode + tools_condition
    build_prebuilt_agent  the same agent from create_agent

Then:

    uv run pytest exercises/06-prebuilt-agent -v
"""
