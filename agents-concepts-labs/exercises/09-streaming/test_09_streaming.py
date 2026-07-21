"""Spec for exercise 09. Do not edit — make it pass."""

from __future__ import annotations

import pytest
from langchain_core.messages import AIMessage, AIMessageChunk

from graph import (
    build_graph,
    stream_answer_tokens,
    stream_node_updates,
    stream_progress,
    stream_state_sizes,
)
from labgraph import ScriptedChatModel, llm_configured, tool_call


def agent(*responses):
    return build_graph(ScriptedChatModel(list(responses)))


def calculator_agent():
    """One tool round trip, then a four-word answer."""
    return agent(tool_call("add", {"a": 2, "b": 3}), AIMessage("The answer is 5"))


# --- stream_mode="updates": what each node did ---------------------------


def test_updates_name_each_node_as_it_finishes():
    assert stream_node_updates(calculator_agent(), "2+3?") == ["model", "tools", "model"]


def test_updates_arrive_during_the_run_not_after():
    # The point of streaming: you see step 1 before step 3 has happened.
    stream = calculator_agent().stream({"messages": [("user", "2+3?")]}, stream_mode="updates")

    first = next(stream)

    assert list(first) == ["model"], "the first chunk arrives before the tool has run"


def test_an_update_carries_only_that_nodes_return_value():
    stream = calculator_agent().stream({"messages": [("user", "2+3?")]}, stream_mode="updates")

    tool_update = [chunk for chunk in stream if "tools" in chunk][0]

    assert [m.content for m in tool_update["tools"]["messages"]] == ["5.0"], (
        "an update is the node's update — not the whole state"
    )


# --- stream_mode="values": the whole state, each step --------------------


def test_values_stream_the_full_state_after_every_step():
    assert stream_state_sizes(calculator_agent(), "2+3?") == [1, 2, 3, 4]


def test_the_last_value_is_what_invoke_would_have_returned():
    graph = calculator_agent()

    states = list(graph.stream({"messages": [("user", "2+3?")]}, stream_mode="values"))

    assert states[-1]["messages"][-1].content == "The answer is 5"


# --- stream_mode="messages": tokens --------------------------------------


def test_the_answer_arrives_in_more_than_one_piece():
    tokens = stream_answer_tokens(calculator_agent(), "2+3?")

    assert len(tokens) > 1, "this is token streaming, not one chunk at the end"
    assert "".join(tokens) == "The answer is 5", "the pieces must rejoin exactly"


def test_tool_output_is_not_mistaken_for_model_tokens():
    # "messages" mode emits ToolMessages too. Rendering "5.0" into the chat
    # bubble is the classic bug.
    tokens = stream_answer_tokens(calculator_agent(), "2+3?")

    assert "5.0" not in tokens


def test_empty_tool_call_chunks_are_not_rendered():
    # The model's first turn is a tool call: an AIMessageChunk with no content.
    # Streamed naively it shows up as a flicker of nothing.
    assert all(token for token in stream_answer_tokens(calculator_agent(), "2+3?"))


def test_messages_mode_reports_which_node_produced_each_chunk():
    graph = calculator_agent()

    pairs = list(graph.stream({"messages": [("user", "2+3?")]}, stream_mode="messages"))

    nodes = {metadata.get("langgraph_node") for _, metadata in pairs}
    assert nodes == {"model", "tools"}
    assert any(isinstance(chunk, AIMessageChunk) for chunk, _ in pairs)


# --- stream_mode="custom": your own progress events ----------------------


def test_custom_events_come_from_get_stream_writer():
    events = stream_progress(calculator_agent(), "2+3?")

    assert events == [
        {"stage": "thinking", "messages_so_far": 1},
        {"stage": "thinking", "messages_so_far": 3},
    ], "one event per model turn, showing the conversation growing"


def test_custom_events_are_absent_from_the_other_modes():
    graph = calculator_agent()

    states = list(graph.stream({"messages": [("user", "2+3?")]}, stream_mode="values"))

    assert all("stage" not in state for state in states), (
        "custom events are a side channel — they never touch state"
    )


# --- several modes at once -----------------------------------------------


def test_you_can_subscribe_to_several_modes():
    graph = calculator_agent()

    tagged = list(
        graph.stream({"messages": [("user", "2+3?")]}, stream_mode=["updates", "custom"])
    )

    modes = [mode for mode, _ in tagged]
    assert set(modes) == {"updates", "custom"}, "chunks are tagged with their mode"
    assert modes[0] == "custom", "the node emits progress before it returns an update"


# --- streaming does not change the answer --------------------------------


def test_streaming_and_invoking_produce_the_same_final_state():
    streamed = list(
        calculator_agent().stream({"messages": [("user", "2+3?")]}, stream_mode="values")
    )[-1]
    invoked = calculator_agent().invoke({"messages": [("user", "2+3?")]})

    assert [m.content for m in streamed["messages"]] == [m.content for m in invoked["messages"]]


# --- optional: against a real provider ------------------------------------


@pytest.mark.skipif(not llm_configured(), reason="no LLM configured in .env")
def test_live_model_streams_tokens():
    from labgraph import get_chat_model

    tokens = stream_answer_tokens(build_graph(get_chat_model()), "Write one short sentence.")

    assert len(tokens) > 1
    assert "".join(tokens).strip()
