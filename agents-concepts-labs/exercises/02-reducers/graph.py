"""Exercise 02 — Reducers: controlling how updates merge.

In exercise 01 every node owned its own key, so "the update overwrites the old
value" was fine. The moment two nodes write the *same* key — or one node writes
it repeatedly in a loop — you have to say what merging means. That is a reducer.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/02-reducers/README.md

The spec imports these names — they must exist with these exact spellings:

    keep_max       a custom reducer: (current, update) -> larger
    RunState       state schema with four different merge behaviours
    stage_one      node writing all four keys
    stage_two      node writing all four keys again
    build_graph    returns the compiled graph

Then:

    uv run pytest exercises/02-reducers -v
"""
