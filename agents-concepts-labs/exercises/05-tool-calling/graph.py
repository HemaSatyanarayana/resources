"""Exercise 05 — Tool calling from scratch: the ReAct loop.

This is the exercise the whole course has been building toward. An agent is not
a special LangGraph object — it is exercise 03's conditional loop with exercise
04's LLM node, and one new question in the router: "did the model ask for a
tool?"

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/05-tool-calling/README.md

The spec imports these names — they must exist with these exact spellings:

    TOOLS            the list of tools the agent may use
    make_model_node  model -> node that calls the tool-bound model
    run_tools        node: execute the requested tool calls
    should_continue  the router — "tools" or END
    build_graph      model -> compiled graph

Then:

    uv run pytest exercises/05-tool-calling -v
"""
