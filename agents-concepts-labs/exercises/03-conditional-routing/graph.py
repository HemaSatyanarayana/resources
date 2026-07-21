"""Exercise 03 — Conditional edges: branching and looping.

Straight lines are not agents. An agent decides what to do next, and in
LangGraph that decision is a *conditional edge*: a function that looks at state
and returns the name of the next node.

You are building a draft/review loop — the skeleton of every "critique and
retry" agent.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/03-conditional-routing/README.md

The spec imports these names — they must exist with these exact spellings:

    MAX_REVISIONS       the budget: 3
    DraftState          the state schema
    write_draft         node: first draft
    revise_draft        node: improve the draft, bump the counter
    review              node: record that review happened
    route_after_review  the conditional edge — loop or stop
    build_graph         returns the compiled graph

Then:

    uv run pytest exercises/03-conditional-routing -v
"""
