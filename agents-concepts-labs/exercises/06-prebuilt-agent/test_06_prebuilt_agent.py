"""Spec for exercise 06. Do not edit — make it pass."""

from __future__ import annotations

import pytest
from langchain_core.messages import AIMessage, SystemMessage, ToolMessage
from langgraph.prebuilt import tools_condition

from graph import SYSTEM_PROMPT, TOOLS, build_graph, build_prebuilt_agent
from labgraph import ScriptedChatModel, llm_configured, tool_call, tool_calls


def script(*responses):
    return ScriptedChatModel(list(responses))


# --- tools_condition is exercise 05's router ------------------------------


def test_tools_condition_is_the_router_you_already_wrote():
    assert tools_condition({"messages": [tool_call("add", {"a": 1, "b": 2})]}) == "tools"
    assert tools_condition({"messages": [AIMessage("done")]}) == "__end__"


# --- the loop still works -------------------------------------------------


def test_graph_runs_a_react_cycle():
    model = script(tool_call("add", {"a": 2, "b": 3}), AIMessage("2 + 3 is 5."))

    result = build_graph(model).invoke({"messages": [("user", "what is 2+3?")]})

    assert [m.type for m in result["messages"]] == ["human", "ai", "tool", "ai"]
    assert result["messages"][-1].content == "2 + 3 is 5."


def test_graph_answers_without_tools():
    model = script(AIMessage("Paris."))

    result = build_graph(model).invoke({"messages": [("user", "capital of France?")]})

    assert result["messages"][-1].content == "Paris."
    assert model.call_count == 1


def test_tool_node_handles_parallel_calls():
    model = script(
        tool_calls(("add", {"a": 1, "b": 1}), ("multiply", {"a": 3, "b": 4})),
        AIMessage("2 and 12."),
    )

    result = build_graph(model).invoke({"messages": [("user", "compute both")]})

    tool_results = [m.content for m in result["messages"] if isinstance(m, ToolMessage)]
    assert tool_results == ["2.0", "12.0"]


# --- what ToolNode gives you for free -------------------------------------


def test_tool_node_turns_an_unknown_tool_into_an_observation():
    # You hand-wrote this in exercise 05. ToolNode does it by default.
    model = script(tool_call("launch_missiles", {}), AIMessage("I cannot do that."))

    result = build_graph(model).invoke({"messages": [("user", "do it")]})

    error = next(m for m in result["messages"] if isinstance(m, ToolMessage))
    assert "error" in error.content.lower()
    assert "launch_missiles" in error.content
    assert result["messages"][-1].content == "I cannot do that.", (
        "the run survives a hallucinated tool and the model gets another turn"
    )


def test_tool_node_turns_bad_arguments_into_an_observation():
    # This one you did NOT write in 05 — a raise from inside the tool.
    model = script(tool_call("add", {"a": "two", "b": 1}), AIMessage("Let me retry."))

    result = build_graph(model).invoke({"messages": [("user", "add two and one")]})

    error = next(m for m in result["messages"] if isinstance(m, ToolMessage))
    assert "error" in error.content.lower()
    assert result["messages"][-1].content == "Let me retry."


# --- the system prompt ----------------------------------------------------


def test_system_prompt_leads_every_model_call():
    model = script(tool_call("add", {"a": 1, "b": 1}), AIMessage("Two."))

    build_graph(model).invoke({"messages": [("user", "1+1?")]})

    for i, call in enumerate(model.calls):
        assert isinstance(call[0], SystemMessage), f"model call {i} had no system message"
        assert call[0].content == SYSTEM_PROMPT


def test_system_prompt_is_not_stored_in_state():
    # Prepend it at call time. Putting it in `messages` means it is checkpointed,
    # re-sent, and duplicated once per turn — and you can never change it.
    model = script(AIMessage("ok"))

    result = build_graph(model).invoke({"messages": [("user", "hi")]})

    assert not any(isinstance(m, SystemMessage) for m in result["messages"])
    assert [m.type for m in result["messages"]] == ["human", "ai"]


def test_system_prompt_is_swappable_per_graph():
    model = script(AIMessage("arrr"))

    build_graph(model, system_prompt="You are a pirate.").invoke(
        {"messages": [("user", "hi")]}
    )

    assert model.calls[0][0].content == "You are a pirate."


# --- the prebuilt ---------------------------------------------------------


def test_prebuilt_agent_runs_the_same_cycle():
    model = script(tool_call("multiply", {"a": 6, "b": 7}), AIMessage("42."))

    result = build_prebuilt_agent(model).invoke({"messages": [("user", "6*7?")]})

    assert [m.type for m in result["messages"]] == ["human", "ai", "tool", "ai"]
    assert result["messages"][-1].content == "42."


def test_prebuilt_agent_gets_the_same_tools_and_system_prompt():
    model = script(AIMessage("hello"))

    build_prebuilt_agent(model).invoke({"messages": [("user", "hi")]})

    assert {t.name for t in model.bound_tools} == {t.name for t in TOOLS}
    assert model.calls[0][0].content == SYSTEM_PROMPT


def test_handwritten_and_prebuilt_agree_message_for_message():
    """The payoff: your graph and the prebuilt are the same machine."""
    responses = [tool_call("search", {"query": "checkpointer"}), AIMessage("It persists state.")]
    question = {"messages": [("user", "what is a checkpointer?")]}

    mine = build_graph(script(*responses)).invoke(question)
    theirs = build_prebuilt_agent(script(*responses)).invoke(question)

    assert [m.type for m in mine["messages"]] == [m.type for m in theirs["messages"]]
    assert [m.content for m in mine["messages"]] == [m.content for m in theirs["messages"]]


def test_both_expose_a_model_node_and_a_tools_node():
    for graph in (build_graph(script()), build_prebuilt_agent(script())):
        assert {"model", "tools"} <= set(graph.get_graph().nodes)


# --- optional: against a real provider ------------------------------------


@pytest.mark.skipif(not llm_configured(), reason="no LLM configured in .env")
def test_live_prebuilt_agent_uses_a_tool():
    from labgraph import get_chat_model

    result = build_prebuilt_agent(get_chat_model()).invoke(
        {"messages": [("user", "What is 128 times 4?")]}
    )

    assert any(isinstance(m, ToolMessage) for m in result["messages"])
    assert "512" in result["messages"][-1].content
