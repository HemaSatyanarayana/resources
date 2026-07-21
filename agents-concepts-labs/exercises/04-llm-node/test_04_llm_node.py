"""Spec for exercise 04. Do not edit — make it pass.

Every test here drives the graph with `ScriptedChatModel`, so the suite is
deterministic and needs no API key. That is the payoff of injecting the model.
"""

from __future__ import annotations

import pytest
from langchain_core.messages import AIMessage, SystemMessage

from graph import (
    CLASSIFIER_PROMPT,
    PERSONAS,
    build_graph,
    make_classifier,
    make_responder,
    route_by_category,
)
from labgraph import ScriptedChatModel, llm_configured


# --- the classifier node --------------------------------------------------


def test_classifier_prompts_with_the_system_instruction_and_conversation():
    model = ScriptedChatModel(["billing"])

    make_classifier(model)({"messages": [("user", "I was charged twice")]})

    prompt = model.last_prompt()
    assert isinstance(prompt[0], SystemMessage)
    assert prompt[0].content == CLASSIFIER_PROMPT
    assert "charged twice" in str(prompt[-1].content)


def test_classifier_writes_the_category_and_nothing_else():
    model = ScriptedChatModel(["technical"])

    update = make_classifier(model)({"messages": [("user", "it crashes")]})

    assert update == {"category": "technical"}, (
        "the classifier's reply is scaffolding — it should not land in messages"
    )


@pytest.mark.parametrize(
    "raw, expected",
    [
        ("billing", "billing"),
        ("Billing.", "billing"),
        ("  TECHNICAL  ", "technical"),
        ("General!", "general"),
    ],
)
def test_classifier_normalises_model_output(raw, expected):
    update = make_classifier(ScriptedChatModel([raw]))({"messages": [("user", "hi")]})

    assert update["category"] == expected


def test_classifier_falls_back_to_general_on_unexpected_output():
    # Models return surprises. "I think this is a billing issue!" is not a
    # category, and the graph must not crash or route somewhere that does not exist.
    model = ScriptedChatModel(["I think this is a billing issue!"])

    update = make_classifier(model)({"messages": [("user", "hi")]})

    assert update["category"] == "general"


# --- the responder node ---------------------------------------------------


def test_responder_uses_its_persona_and_appends_the_reply():
    model = ScriptedChatModel([AIMessage("Your invoice is attached.")])

    update = make_responder(model, "billing")({"messages": [("user", "invoice?")]})

    assert model.last_prompt()[0].content == PERSONAS["billing"]
    assert [m.content for m in update["messages"]] == ["Your invoice is attached."]


def test_each_persona_is_distinct():
    for category in PERSONAS:
        model = ScriptedChatModel(["ok"])
        make_responder(model, category)({"messages": [("user", "hi")]})
        assert model.last_prompt()[0].content == PERSONAS[category]


# --- routing --------------------------------------------------------------


def test_router_returns_the_category():
    assert route_by_category({"category": "technical", "messages": []}) == "technical"


# --- end to end -----------------------------------------------------------


def test_graph_classifies_then_answers_in_that_persona():
    # Script: first call is the classifier, second is the responder.
    model = ScriptedChatModel(["billing", AIMessage("Refund issued.")])

    result = build_graph(model).invoke({"messages": [("user", "I was charged twice")]})

    assert result["category"] == "billing"
    assert result["messages"][-1].content == "Refund issued."
    assert model.call_count == 2, "exactly one classify call and one respond call"


def test_graph_keeps_the_original_user_message():
    model = ScriptedChatModel(["general", AIMessage("Hello!")])

    result = build_graph(model).invoke({"messages": [("user", "hi there")]})

    assert [m.content for m in result["messages"]] == ["hi there", "Hello!"]


def test_graph_routes_to_the_technical_persona():
    model = ScriptedChatModel(["technical", AIMessage("Clear your cache.")])
    graph = build_graph(model)

    graph.invoke({"messages": [("user", "the app crashes on launch")]})

    # The second model call is the responder — check which persona it got.
    assert model.calls[1][0].content == PERSONAS["technical"]


def test_only_one_responder_runs():
    # A conditional edge picks ONE branch. If all three responders ran, the
    # script would be exhausted and the model would raise.
    model = ScriptedChatModel(["general", AIMessage("Sure thing.")])

    result = build_graph(model).invoke({"messages": [("user", "hello")]})

    assert len(result["messages"]) == 2


def test_graph_exposes_all_four_nodes():
    nodes = set(build_graph(ScriptedChatModel([])).get_graph().nodes)

    assert {"classify", "billing", "technical", "general"} <= nodes


# --- optional: against a real provider ------------------------------------


@pytest.mark.skipif(not llm_configured(), reason="no LLM configured in .env")
def test_live_model_classifies_a_billing_request():
    """Only runs when .env is filled in. Everything above stays offline."""
    from labgraph import get_chat_model

    result = build_graph(get_chat_model()).invoke(
        {"messages": [("user", "My card was charged twice for the same invoice")]}
    )

    assert result["category"] in ("billing", "general")
    assert result["messages"][-1].content
