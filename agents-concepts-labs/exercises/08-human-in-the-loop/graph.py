"""Exercise 08 — Human in the loop: interrupt() and Command(resume=...).

Your agent will call any tool the model asks for. That is fine when the tools
read. It is not fine when a tool spends money. You are building an approval
gate: `search` runs freely, `refund` needs a human.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/08-human-in-the-loop/README.md

The spec imports these names — they must exist with these exact spellings:

    LEDGER                stand-in for a payments API; proves what really ran
    refund                the @tool that spends money
    TOOLS                 the list of tools the agent may use
    SENSITIVE_TOOLS       the names that require approval
    ApprovalState         MessagesState plus a `decision` key
    make_model_node       model -> node that calls the tool-bound model
    review_tool_calls     node: interrupt for a human when needed
    route_after_model     the router into the gate
    route_after_approval  the router out of the gate
    build_graph           (model, checkpointer=None) -> compiled graph

Then:

    uv run pytest exercises/08-human-in-the-loop -v
"""
