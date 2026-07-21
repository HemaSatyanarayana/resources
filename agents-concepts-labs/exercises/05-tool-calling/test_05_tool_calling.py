"""Spec for exercise 05. Do not edit — make it pass."""

from __future__ import annotations

import pytest
from langchain_core.messages import AIMessage, ToolMessage
from langgraph.errors import GraphRecursionError

from graph import TOOLS, build_graph, make_model_node, run_tools, should_continue
from labgraph import ScriptedChatModel, llm_configured, tool_call, tool_calls


# --- the model node -------------------------------------------------------


def test_model_node_binds_the_tools():
    model = ScriptedChatModel(["hi"])

    make_model_node(model)({"messages": [("user", "hi")]})

    assert {t.name for t in model.bound_tools} == {t.name for t in TOOLS}, (
        "the model cannot request a tool it was never told about"
    )


def test_model_node_appends_the_reply():
    model = ScriptedChatModel([AIMessage("hello")])

    update = make_model_node(model)({"messages": [("user", "hi")]})

    assert [m.content for m in update["messages"]] == ["hello"]


def test_tools_are_bound_once_not_per_call():
    # bind_tools() is setup, not work. Doing it inside the node re-does it on
    # every turn of the loop.
    model = ScriptedChatModel(["a", "b"], loop=True)
    node = make_model_node(model)

    node({"messages": [("user", "1")]})
    bound_after_first = model.bound_tools

    node({"messages": [("user", "2")]})

    assert model.bound_tools is bound_after_first


# --- the tool node --------------------------------------------------------


def test_run_tools_executes_the_call_and_returns_a_tool_message():
    state = {"messages": [tool_call("add", {"a": 2, "b": 3})]}

    update = run_tools(state)

    (message,) = update["messages"]
    assert isinstance(message, ToolMessage)
    assert message.content == "5.0"
    assert message.tool_call_id == "call_add", (
        "the id links the result back to the request — providers reject a "
        "mismatch"
    )


def test_run_tools_handles_parallel_calls_in_order():
    state = {"messages": [tool_calls(("add", {"a": 1, "b": 1}), ("multiply", {"a": 3, "b": 4}))]}

    update = run_tools(state)

    assert [m.content for m in update["messages"]] == ["2.0", "12.0"]
    assert [m.tool_call_id for m in update["messages"]] == ["call_add_0", "call_multiply_1"]


def test_run_tools_reads_the_last_message_only():
    state = {
        "messages": [
            ("user", "what is 2+3?"),
            tool_call("add", {"a": 2, "b": 3}),
        ]
    }

    assert run_tools(state)["messages"][0].content == "5.0"


def test_run_tools_reports_an_unknown_tool_instead_of_crashing():
    # A model can hallucinate a tool name. Raising kills the run; answering with
    # an error message lets the model see what it did and correct itself.
    state = {"messages": [tool_call("delete_everything", {})]}

    update = run_tools(state)

    (message,) = update["messages"]
    assert isinstance(message, ToolMessage)
    assert "unknown tool" in message.content.lower()
    assert "delete_everything" in message.content


# --- the router -----------------------------------------------------------


def test_router_goes_to_tools_when_the_model_asked_for_one():
    assert should_continue({"messages": [tool_call("search", {"query": "x"})]}) == "tools"


def test_router_ends_on_a_plain_reply():
    assert should_continue({"messages": [AIMessage("the answer is 5")]}) == "__end__"


def test_router_ends_on_a_tool_result_the_model_did_not_follow_up():
    # Defensive: the last message is not always an AIMessage.
    state = {"messages": [ToolMessage(content="5", tool_call_id="x", name="add")]}

    assert should_continue(state) == "__end__"


# --- the loop -------------------------------------------------------------


def test_graph_answers_directly_without_touching_a_tool():
    model = ScriptedChatModel([AIMessage("Paris.")])

    result = build_graph(model).invoke({"messages": [("user", "capital of France?")]})

    assert result["messages"][-1].content == "Paris."
    assert model.call_count == 1
    assert not any(isinstance(m, ToolMessage) for m in result["messages"])


def test_graph_runs_one_full_react_cycle():
    model = ScriptedChatModel(
        [tool_call("add", {"a": 2, "b": 3}), AIMessage("2 + 3 is 5.")]
    )

    result = build_graph(model).invoke({"messages": [("user", "what is 2+3?")]})

    kinds = [m.type for m in result["messages"]]
    assert kinds == ["human", "ai", "tool", "ai"], (
        "user asks -> model requests -> tool answers -> model summarises"
    )
    assert result["messages"][-1].content == "2 + 3 is 5."
    assert model.call_count == 2


def test_the_model_sees_the_tool_result_on_the_second_turn():
    # This is the whole point of the loop: the result must get back into the
    # prompt, or the model is answering from nothing.
    model = ScriptedChatModel(
        [tool_call("search", {"query": "what is a reducer"}), AIMessage("Done.")]
    )

    build_graph(model).invoke({"messages": [("user", "what is a reducer?")]})

    assert "reducer decides how" in model.last_prompt_text()


def test_graph_loops_as_many_times_as_the_model_asks():
    model = ScriptedChatModel(
        [
            tool_call("add", {"a": 1, "b": 1}),
            tool_call("multiply", {"a": 2, "b": 5}),
            AIMessage("Ten."),
        ]
    )

    result = build_graph(model).invoke({"messages": [("user", "compute something")]})

    assert [m.type for m in result["messages"]] == [
        "human", "ai", "tool", "ai", "tool", "ai",
    ]
    assert model.call_count == 3


def test_a_model_that_never_stops_hits_the_recursion_limit():
    # The loop has no budget of its own — `recursion_limit` is the backstop.
    # Each scripted call needs its own id: `add_messages` merges by id, so
    # replaying one identical message would overwrite rather than append.
    model = ScriptedChatModel(
        [tool_call("add", {"a": 1, "b": 1}, id=f"call_{i}") for i in range(20)]
    )

    with pytest.raises(GraphRecursionError):
        build_graph(model).invoke(
            {"messages": [("user", "go forever")]}, {"recursion_limit": 8}
        )


def test_graph_has_a_model_node_and_a_tools_node():
    nodes = set(build_graph(ScriptedChatModel([])).get_graph().nodes)

    assert {"model", "tools"} <= nodes


# --- optional: against a real provider ------------------------------------


@pytest.mark.skipif(not llm_configured(), reason="no LLM configured in .env")
def test_live_model_uses_the_calculator():
    from labgraph import get_chat_model

    result = build_graph(get_chat_model()).invoke(
        {"messages": [("user", "Use your tools: what is 128 multiplied by 4?")]}
    )

    assert any(isinstance(m, ToolMessage) for m in result["messages"])
    assert "512" in result["messages"][-1].content
