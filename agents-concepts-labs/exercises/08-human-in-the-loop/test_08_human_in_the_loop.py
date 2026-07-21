"""Spec for exercise 08. Do not edit — make it pass."""

from __future__ import annotations

import pytest
from langchain_core.messages import AIMessage, ToolMessage
from langgraph.checkpoint.memory import InMemorySaver
from langgraph.types import Command

from graph import (
    LEDGER,
    build_graph,
    route_after_approval,
    route_after_model,
)
from labgraph import ScriptedChatModel, llm_configured, tool_call


@pytest.fixture(autouse=True)
def clean_ledger():
    LEDGER.clear()
    yield
    LEDGER.clear()


def thread(name="t"):
    return {"configurable": {"thread_id": name}}


def agent(*responses):
    """A graph with memory — an interrupt cannot exist without a checkpointer."""
    return build_graph(ScriptedChatModel(list(responses)), InMemorySaver())


def refund_call(customer="ana", amount=20.0):
    return tool_call("refund", {"customer": customer, "amount": amount})


# --- routing --------------------------------------------------------------


def test_tool_calls_go_through_approval_not_straight_to_tools():
    assert route_after_model({"messages": [refund_call()]}) == "approval"


def test_a_plain_answer_ends_the_run():
    assert route_after_model({"messages": [AIMessage("all done")]}) == "__end__"


def test_approved_runs_the_tools_denied_goes_back_to_the_model():
    assert route_after_approval({"messages": [], "decision": "approved"}) == "tools"
    assert route_after_approval({"messages": [], "decision": "denied"}) == "model"


# --- safe tools are not gated --------------------------------------------


def test_a_read_only_tool_never_pauses():
    graph = agent(tool_call("search", {"query": "reducer"}), AIMessage("It merges updates."))

    result = graph.invoke({"messages": [("user", "what is a reducer?")]}, thread())

    assert "__interrupt__" not in result, "only SENSITIVE_TOOLS need a human"
    assert result["messages"][-1].content == "It merges updates."


# --- the pause ------------------------------------------------------------


def test_a_sensitive_tool_pauses_the_run():
    graph = agent(refund_call(), AIMessage("Refunded."))

    result = graph.invoke({"messages": [("user", "refund ana $20")]}, thread())

    assert "__interrupt__" in result, "the run should stop before spending money"
    assert LEDGER == [], "the tool must NOT have run yet"


def test_the_pause_carries_a_payload_a_human_can_act_on():
    graph = agent(refund_call("ana", 20.0), AIMessage("Refunded."))

    result = graph.invoke({"messages": [("user", "refund ana $20")]}, thread())

    payload = result["__interrupt__"][0].value
    assert payload["tool_calls"] == [
        {"name": "refund", "args": {"customer": "ana", "amount": 20.0}}
    ], "a reviewer needs to see exactly what they are approving"


def test_a_paused_graph_reports_what_it_is_waiting_on():
    graph = agent(refund_call(), AIMessage("Refunded."))
    graph.invoke({"messages": [("user", "refund ana $20")]}, thread())

    snapshot = graph.get_state(thread())

    assert snapshot.next == ("approval",), (
        "exercise 07's `.next` is empty for a finished run — this run is not finished"
    )


def test_the_pause_survives_being_left_alone():
    # Nothing is held in memory between the two calls but the checkpoint. A real
    # approval might arrive days later, from a different process.
    graph = agent(refund_call(), AIMessage("Refunded."))
    graph.invoke({"messages": [("user", "refund ana $20")]}, thread())

    assert graph.get_state(thread()).values["messages"][-1].tool_calls


# --- approving ------------------------------------------------------------


def test_resuming_with_approval_runs_the_tool():
    graph = agent(refund_call("ana", 20.0), AIMessage("Refunded $20 to ana."))
    graph.invoke({"messages": [("user", "refund ana $20")]}, thread())

    result = graph.invoke(Command(resume=True), thread())

    assert LEDGER == [{"customer": "ana", "amount": 20.0}]
    assert result["messages"][-1].content == "Refunded $20 to ana."
    assert [m.type for m in result["messages"]] == ["human", "ai", "tool", "ai"]


def test_resuming_does_not_replay_the_model_call():
    # The graph resumes at the interrupt, not from the top. The first model call
    # already happened and is not repeated.
    model = ScriptedChatModel([refund_call(), AIMessage("Done.")])
    graph = build_graph(model, InMemorySaver())

    graph.invoke({"messages": [("user", "refund ana $20")]}, thread())
    assert model.call_count == 1

    graph.invoke(Command(resume=True), thread())
    assert model.call_count == 2, "one more call to summarise — not a restart"


# --- denying --------------------------------------------------------------


def test_resuming_with_a_denial_never_runs_the_tool():
    graph = agent(refund_call(), AIMessage("Understood, I won't issue the refund."))
    graph.invoke({"messages": [("user", "refund ana $20")]}, thread())

    result = graph.invoke(Command(resume=False), thread())

    assert LEDGER == [], "denied means it never ran"
    assert result["messages"][-1].content == "Understood, I won't issue the refund."


def test_a_denial_is_reported_back_to_the_model_as_a_tool_result():
    # Every tool call needs a matching ToolMessage — providers reject a dangling
    # call. And the model has to learn what happened, or it will just try again.
    graph = agent(refund_call(), AIMessage("Understood."))
    graph.invoke({"messages": [("user", "refund ana $20")]}, thread())

    result = graph.invoke(Command(resume=False), thread())

    denial = next(m for m in result["messages"] if isinstance(m, ToolMessage))
    assert "denied" in denial.content.lower()
    assert denial.tool_call_id == "call_refund", "it answers the call it refused"


def test_the_model_sees_the_denial_on_its_next_turn():
    model = ScriptedChatModel([refund_call(), AIMessage("Understood.")])
    graph = build_graph(model, InMemorySaver())
    graph.invoke({"messages": [("user", "refund ana $20")]}, thread())

    graph.invoke(Command(resume=False), thread())

    assert "denied" in model.last_prompt_text()


# --- threads are still threads -------------------------------------------


def test_two_reviews_can_be_pending_at_once():
    graph = agent(refund_call("ana", 20.0), refund_call("bo", 5.0), AIMessage("Done."))

    graph.invoke({"messages": [("user", "refund ana $20")]}, thread("ana"))
    graph.invoke({"messages": [("user", "refund bo $5")]}, thread("bo"))

    assert graph.get_state(thread("ana")).next == ("approval",)
    assert graph.get_state(thread("bo")).next == ("approval",)

    graph.invoke(Command(resume=True), thread("bo"))

    assert LEDGER == [{"customer": "bo", "amount": 5.0}], "approving one is not approving both"
    assert graph.get_state(thread("ana")).next == ("approval",)


# --- optional: against a real provider ------------------------------------


@pytest.mark.skipif(not llm_configured(), reason="no LLM configured in .env")
def test_live_model_gets_gated_on_a_refund():
    from labgraph import get_chat_model

    graph = build_graph(get_chat_model(), InMemorySaver())

    result = graph.invoke(
        {"messages": [("user", "Please refund $20 to customer ana.")]}, thread("live")
    )

    assert "__interrupt__" in result
    assert LEDGER == []
