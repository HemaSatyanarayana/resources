"""Exercise 04 — Putting a language model inside a node.

A node that calls an LLM is still just `state -> update`. What changes is that
the update is now non-deterministic — which forces two habits: inject the model
rather than constructing it, and parse the model's prose into state defensively.

You are building a support-desk router.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/04-llm-node/README.md

The spec imports these names — they must exist with these exact spellings:

    CATEGORIES         ("billing", "technical", "general")
    CLASSIFIER_PROMPT  the system prompt for the classifier
    PERSONAS           category -> system prompt for that responder
    SupportState       MessagesState plus a `category` key
    make_classifier    model -> node that writes `category`
    make_responder     (model, category) -> node that answers in persona
    route_by_category  the conditional edge
    build_graph        model -> compiled graph

Then:

    uv run pytest exercises/04-llm-node -v
"""
