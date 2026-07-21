"""Spec for exercise 10. Do not edit — make it pass."""

from __future__ import annotations

import operator
from typing import Annotated, TypedDict

import pytest
from langchain_core.messages import AIMessage
from langgraph.errors import InvalidUpdateError
from langgraph.graph import END, START, StateGraph
from langgraph.types import Send

from graph import (
    RESEARCHER_PROMPT,
    build_graph,
    fan_out,
    make_planner,
    make_researcher,
    parse_subtopics,
    write_report,
)
from labgraph import ScriptedChatModel, llm_configured


def start(topic="AI"):
    return {"topic": topic, "subtopics": [], "notes": [], "report": ""}


def pipeline(plan: str, *notes: str):
    """A model that first returns a plan, then one note per worker."""
    return ScriptedChatModel([AIMessage(plan), *(AIMessage(n) for n in notes)])


# --- parsing --------------------------------------------------------------


@pytest.mark.parametrize(
    "raw, expected",
    [
        ("history, ethics, jobs", ["history", "ethics", "jobs"]),
        ("  history ,ethics  ", ["history", "ethics"]),
        ("solo", ["solo"]),
        ("a,,b", ["a", "b"]),
        ("   ", []),
    ],
)
def test_parse_subtopics_is_defensive(raw, expected):
    assert parse_subtopics(raw) == expected


# --- planning -------------------------------------------------------------


def test_planner_writes_subtopics_only():
    model = pipeline("history, ethics, jobs")

    update = make_planner(model)(start("AI"))

    assert update == {"subtopics": ["history", "ethics", "jobs"]}
    assert model.call_count == 1


# --- the fan-out ----------------------------------------------------------


def test_fan_out_returns_one_send_per_subtopic():
    sends = fan_out({**start(), "subtopics": ["a", "b", "c"]})

    assert all(isinstance(s, Send) for s in sends)
    assert [s.node for s in sends] == ["research", "research", "research"], (
        "same node, three times — a Send is a task, not an edge"
    )


def test_each_send_carries_only_its_own_work():
    sends = fan_out({**start(), "subtopics": ["a", "b"]})

    assert [s.arg for s in sends] == [{"subtopic": "a"}, {"subtopic": "b"}], (
        "a worker gets the payload you hand it, not the parent state"
    )


def test_the_number_of_workers_is_decided_at_runtime():
    assert len(fan_out({**start(), "subtopics": list("abcde")})) == 5
    assert len(fan_out({**start(), "subtopics": ["only"]})) == 1


def test_fan_out_skips_straight_to_write_when_there_is_nothing_to_do():
    assert fan_out({**start(), "subtopics": []}) == "write", (
        "a conditional edge may return Sends OR a plain node name"
    )


# --- the worker -----------------------------------------------------------


def test_worker_sees_its_subtopic_and_nothing_else():
    model = ScriptedChatModel([AIMessage("Ethics matters.")])

    update = make_researcher(model)({"subtopic": "ethics"})

    prompt = [str(m.content) for m in model.last_prompt()]
    assert prompt[0] == RESEARCHER_PROMPT
    assert "ethics" in prompt[-1]
    assert "AI" not in prompt[-1], "the worker was never given the parent topic"


def test_worker_returns_a_single_element_list():
    # A list, because the reducer concatenates. Returning a bare string here
    # would make `notes` a soup of characters.
    model = ScriptedChatModel([AIMessage("Ethics matters.")])

    update = make_researcher(model)({"subtopic": "ethics"})

    assert update == {"notes": ["ethics: Ethics matters."]}


# --- the reduce step ------------------------------------------------------


def test_notes_from_every_worker_survive():
    graph = build_graph(pipeline("history, ethics, jobs", "N1", "N2", "N3"))

    result = graph.invoke(start("AI"))

    assert result["notes"] == ["history: N1", "ethics: N2", "jobs: N3"], (
        "three concurrent writes to one key, and none of them lost"
    )


def test_one_model_call_per_worker_plus_the_planner():
    model = pipeline("a, b, c, d", "1", "2", "3", "4")

    build_graph(model).invoke(start())

    assert model.call_count == 5


def test_report_gathers_every_note():
    graph = build_graph(pipeline("history, ethics", "N1", "N2"))

    report = graph.invoke(start("AI"))["report"]

    assert "# AI" in report
    assert "history: N1" in report
    assert "ethics: N2" in report


def test_write_report_handles_an_empty_run():
    assert "no findings" in write_report({**start("X"), "notes": []})["report"]


def test_an_empty_plan_runs_no_workers():
    model = pipeline("   ")

    result = build_graph(model).invoke(start("X"))

    assert model.call_count == 1, "no subtopics, no workers"
    assert "no findings" in result["report"]


# --- why the reducer is not optional --------------------------------------


def test_concurrent_writes_without_a_reducer_are_an_error():
    """Self-contained proof — this graph does not use your code.

    Two workers, one plain `str` key, no reducer. LangGraph refuses to guess
    which write wins.
    """

    class Broken(TypedDict):
        result: str  # no Annotated[..., reducer]

    def worker(state):
        return {"result": "x"}

    builder = StateGraph(Broken)
    builder.add_node("worker", worker)
    builder.add_conditional_edges(
        START, lambda s: [Send("worker", {}), Send("worker", {})], ["worker"]
    )
    builder.add_edge("worker", END)

    with pytest.raises(InvalidUpdateError):
        builder.compile().invoke({"result": ""})


def test_the_same_graph_works_once_a_reducer_is_added():
    class Fixed(TypedDict):
        result: Annotated[list[str], operator.add]

    def worker(state):
        return {"result": ["x"]}

    builder = StateGraph(Fixed)
    builder.add_node("worker", worker)
    builder.add_conditional_edges(
        START, lambda s: [Send("worker", {}), Send("worker", {})], ["worker"]
    )
    builder.add_edge("worker", END)

    assert builder.compile().invoke({"result": []})["result"] == ["x", "x"]


# --- it is one superstep --------------------------------------------------


def test_the_workers_run_in_one_superstep_then_write_runs_once():
    graph = build_graph(pipeline("a, b", "NA", "NB"))

    chunks = list(graph.stream(start("T"), stream_mode="updates"))
    order = [name for chunk in chunks for name in chunk]

    assert order == ["plan", "research", "research", "write"]
    assert order.count("write") == 1, (
        "`write` waits for BOTH workers — it does not run once per worker"
    )


# --- optional: against a real provider ------------------------------------


@pytest.mark.skipif(not llm_configured(), reason="no LLM configured in .env")
def test_live_model_researches_in_parallel():
    from labgraph import get_chat_model

    result = build_graph(get_chat_model()).invoke(start("the Apollo program"))

    assert len(result["subtopics"]) >= 2
    assert len(result["notes"]) == len(result["subtopics"])
    assert result["report"].startswith("# the Apollo program")
