"""Spec for exercise 07. Do not edit — make it pass."""

from __future__ import annotations

import pytest
from langchain_core.messages import AIMessage
from langgraph.checkpoint.memory import InMemorySaver

from graph import (
    build_chat_graph,
    checkpoint_configs,
    conversation,
    say,
    thread,
)
from labgraph import ScriptedChatModel, llm_configured


def replies(*texts):
    return ScriptedChatModel([AIMessage(t) for t in texts])


def saved(model):
    """A graph with memory, and the checkpointer holding it."""
    return build_chat_graph(model, InMemorySaver())


# --- the config -----------------------------------------------------------


def test_thread_builds_the_config_langgraph_expects():
    assert thread("abc") == {"configurable": {"thread_id": "abc"}}


# --- without a checkpointer, nothing is remembered ------------------------


def test_a_graph_without_a_checkpointer_forgets_immediately():
    model = replies("Hi Hema!", "I don't know your name.")
    graph = build_chat_graph(model)  # no checkpointer

    say(graph, "a", "my name is Hema")
    say(graph, "a", "what is my name?")

    second_prompt = [m.content for m in model.calls[1]]
    assert "my name is Hema" not in second_prompt, (
        "same thread_id, but with no checkpointer there is nowhere to store it"
    )


def test_a_graph_without_a_checkpointer_still_runs():
    graph = build_chat_graph(replies("hello"))

    assert say(graph, "a", "hi") == "hello"


# --- with a checkpointer, the thread accumulates -------------------------


def test_the_second_turn_sees_the_first():
    model = replies("Hi Hema!", "You are Hema.")
    graph = saved(model)

    say(graph, "a", "my name is Hema")
    say(graph, "a", "what is my name?")

    second_prompt = [m.content for m in model.calls[1]]
    assert "my name is Hema" in second_prompt
    assert "Hi Hema!" in second_prompt, "its own replies are part of the history too"


def test_state_accumulates_across_invocations():
    graph = saved(replies("A1", "A2"))

    say(graph, "a", "Q1")
    say(graph, "a", "Q2")

    assert [m.content for m in conversation(graph, "a")] == ["Q1", "A1", "Q2", "A2"]


def test_each_invocation_only_passes_the_new_message():
    # You send one message; the checkpointer supplies the rest. This is the
    # whole ergonomic win — the caller does not carry the transcript around.
    graph = saved(replies("A1", "A2"))

    say(graph, "a", "Q1")
    result = graph.invoke({"messages": [("user", "Q2")]}, thread("a"))

    assert len(result["messages"]) == 4


# --- threads are isolated -------------------------------------------------


def test_different_threads_do_not_see_each_other():
    model = replies("Hi Hema!", "Hi Sam!")
    graph = saved(model)

    say(graph, "hema", "my name is Hema")
    say(graph, "sam", "my name is Sam")

    assert "my name is Hema" not in [m.content for m in model.calls[1]]
    assert [m.content for m in conversation(graph, "hema")] == ["my name is Hema", "Hi Hema!"]
    assert [m.content for m in conversation(graph, "sam")] == ["my name is Sam", "Hi Sam!"]


def test_an_untouched_thread_is_empty():
    graph = saved(replies("hi"))

    assert conversation(graph, "nobody-has-used-this") == []


def test_a_checkpointed_graph_requires_a_thread_id():
    # Persistence needs a key. Forgetting it is an error, not a silent no-op.
    graph = saved(replies("hi"))

    with pytest.raises(Exception):
        graph.invoke({"messages": [("user", "hi")]})


# --- inspecting state -----------------------------------------------------


def test_get_state_reports_a_finished_run():
    graph = saved(replies("A1"))
    say(graph, "a", "Q1")

    snapshot = graph.get_state(thread("a"))

    assert snapshot.next == (), "nothing is pending — the run completed"
    assert snapshot.config["configurable"]["thread_id"] == "a"
    assert snapshot.config["configurable"]["checkpoint_id"], "state is a point in history"


def test_history_holds_every_step_newest_first():
    graph = saved(replies("A1", "A2"))
    say(graph, "a", "Q1")
    say(graph, "a", "Q2")

    sizes = [len(s.values.get("messages", [])) for s in graph.get_state_history(thread("a"))]

    assert len(sizes) > 2, "a checkpoint is written per step, not per invocation"
    assert sizes == sorted(sizes, reverse=True), "history runs backwards in time"
    assert sizes[0] == 4


def test_every_checkpoint_is_individually_addressable():
    graph = saved(replies("A1", "A2"))
    say(graph, "a", "Q1")
    say(graph, "a", "Q2")

    configs = checkpoint_configs(graph, "a")
    ids = [c["configurable"]["checkpoint_id"] for c in configs]

    assert len(ids) == len(set(ids))
    assert all(c["configurable"]["thread_id"] == "a" for c in configs)


# --- time travel ----------------------------------------------------------


def test_you_can_resume_from_an_earlier_checkpoint():
    """Rewind the conversation and take turn 2 down a different road."""
    graph = saved(replies("A1", "A2", "A3"))
    say(graph, "a", "Q1")
    say(graph, "a", "Q2")

    # End of turn 1: two messages, nothing pending.
    turn_1 = next(
        s
        for s in graph.get_state_history(thread("a"))
        if s.next == () and len(s.values.get("messages", [])) == 2
    )

    result = graph.invoke({"messages": [("user", "Q2-alternate")]}, turn_1.config)

    assert [m.content for m in result["messages"]] == ["Q1", "A1", "Q2-alternate", "A3"], (
        "Q2/A2 never happened on this branch"
    )


# --- optional: against a real provider ------------------------------------


@pytest.mark.skipif(not llm_configured(), reason="no LLM configured in .env")
def test_live_model_remembers_across_turns():
    from labgraph import get_chat_model

    graph = build_chat_graph(get_chat_model(), InMemorySaver())

    say(graph, "live", "My name is Hema. Remember it.")
    answer = say(graph, "live", "What is my name?")

    assert "hema" in answer.lower()
