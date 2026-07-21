"""Spec for exercise 02. Do not edit — make it pass."""

from __future__ import annotations

from typing import get_type_hints

from langchain_core.messages import AIMessage, HumanMessage

from graph import RunState, build_graph, keep_max


# --- the custom reducer in isolation --------------------------------------


def test_keep_max_keeps_the_larger_value():
    assert keep_max(3, 7) == 7
    assert keep_max(7, 3) == 7


def test_keep_max_handles_the_first_write():
    # A reducer is an ordinary function and should defend its own contract.
    # (In this graph LangGraph actually seeds `current` with the annotated
    # type's empty value — 0 for an int — rather than None. See the README.)
    assert keep_max(None, 4) == 4


# --- reducers wired into state --------------------------------------------


def _run(initial=None):
    return build_graph().invoke(
        initial
        or {
            "status": "start",
            "log": ["initial"],
            "messages": [HumanMessage("hi")],
            "high_score": 0,
        }
    )


def test_key_without_a_reducer_is_overwritten():
    # Both nodes wrote `status`; the last writer wins and nothing accumulates.
    assert _run()["status"] == "two"


def test_operator_add_accumulates_the_log():
    log = _run()["log"]

    assert log == ["initial", "entered one", "entered two"], (
        "log should concatenate in execution order, keeping the initial value"
    )


def test_add_messages_appends_to_history():
    messages = _run()["messages"]

    assert [m.content for m in messages] == [
        "hi",
        "hello from one",
        "hello from two",
    ]


def test_custom_reducer_keeps_the_high_score():
    # stage_one wrote 10, stage_two wrote 5. keep_max means 10 survives.
    assert _run()["high_score"] == 10


# --- what makes add_messages special --------------------------------------


def test_add_messages_replaces_a_message_with_the_same_id():
    # This is why chat history uses add_messages rather than operator.add:
    # re-emitting a message with a known id edits it in place instead of
    # duplicating it. That is how streaming updates and message edits work.
    graph = build_graph()
    original = AIMessage(content="first draft", id="msg-1")

    result = graph.invoke(
        {
            "status": "start",
            "log": [],
            "messages": [original, AIMessage(content="corrected", id="msg-1")],
            "high_score": 0,
        }
    )

    by_id = [m for m in result["messages"] if m.id == "msg-1"]
    assert len(by_id) == 1, "same id should replace, not append"
    assert by_id[0].content == "corrected"


def test_add_messages_coerces_tuples_into_message_objects():
    # add_messages also normalises input: ("user", "text") becomes a HumanMessage.
    result = build_graph().invoke(
        {"status": "s", "log": [], "messages": [("user", "typed as a tuple")], "high_score": 0}
    )

    first = result["messages"][0]
    assert isinstance(first, HumanMessage)
    assert first.content == "typed as a tuple"


def test_state_schema_declares_reducers_on_the_right_keys():
    # Reducers live in the *type annotation*, so they are part of the schema.
    # include_extras=True keeps the Annotated[...] wrapper instead of erasing it
    # down to the bare type — that wrapper is where the reducer hides.
    hints = get_type_hints(RunState, include_extras=True)

    assert getattr(hints["status"], "__metadata__", None) is None, (
        "status should have no reducer — it is a last-writer-wins key"
    )
    for key in ("log", "messages", "high_score"):
        assert getattr(hints[key], "__metadata__", None), (
            f"{key} needs an Annotated[...] reducer"
        )
