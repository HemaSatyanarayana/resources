"""Exercise 10 — Parallel fan-out: Send and map-reduce.

Every graph so far had a fixed shape: the nodes existed before the run started.
But "research each of these subtopics" cannot be wired ahead of time — you do
not know there are three of them until the planner says so. `Send` is how a
graph decides its own width at runtime.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/10-parallel-fanout/README.md

The spec imports these names — they must exist with these exact spellings:

    PLANNER_PROMPT     system prompt for the planner
    RESEARCHER_PROMPT  system prompt for one worker
    ResearchState      the parent state (with a reducer on `notes`)
    WorkerState        what ONE worker sees — and it is not ResearchState
    parse_subtopics    defensive parsing of the planner's reply
    make_planner       model -> node that writes `subtopics`
    fan_out            the conditional edge that returns N Sends
    make_researcher    model -> the worker node
    write_report       the reduce step
    build_graph        model -> compiled graph

Then:

    uv run pytest exercises/10-parallel-fanout -v
"""
