"""Exercise 01 — Your first StateGraph.

Build a three-node text pipeline. No LLM involved: the point is to internalise
the execution model before anything non-deterministic enters the picture.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/01-state-graph/README.md

The spec imports these names — they must exist with these exact spellings:

    PipelineState    the state schema
    clean_text       node: normalise `text` into `cleaned`
    count_words      node: count words of `cleaned` into `word_count`
    summarize        node: format `summary`
    build_pipeline   returns the compiled graph

Then:

    uv run pytest exercises/01-state-graph -v
"""
