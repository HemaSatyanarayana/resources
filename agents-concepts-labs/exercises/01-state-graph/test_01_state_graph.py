"""Spec for exercise 01. Do not edit — make it pass."""

from __future__ import annotations

import pytest
from langgraph.graph import END, START

from graph import build_pipeline, clean_text, count_words, summarize


# --- nodes in isolation ---------------------------------------------------
# A node is just a function. You can unit test it without a graph at all, which
# is exactly why LangGraph keeps them plain callables.


def test_clean_text_normalises_case_and_whitespace():
    assert clean_text({"text": "  Hello   WORLD  "})["cleaned"] == "hello world"


def test_clean_text_returns_only_its_own_key():
    update = clean_text({"text": "Anything"})
    assert set(update) == {"cleaned"}, (
        "return a partial update holding just the key this node owns"
    )


def test_count_words_reads_the_cleaned_key():
    assert count_words({"cleaned": "one two three"})["word_count"] == 3


def test_count_words_handles_empty_text():
    assert count_words({"cleaned": ""})["word_count"] == 0


def test_summarize_formats_the_report():
    update = summarize({"cleaned": "hello world", "word_count": 2})
    assert update["summary"] == "2 words: hello world"


# --- the compiled graph ---------------------------------------------------


@pytest.fixture
def graph():
    return build_pipeline()


def test_graph_runs_end_to_end(graph):
    result = graph.invoke({"text": "  LangGraph is  EXPLICIT  "})

    assert result["cleaned"] == "langgraph is explicit"
    assert result["word_count"] == 3
    assert result["summary"] == "3 words: langgraph is explicit"


def test_invoke_returns_the_whole_state_not_just_the_last_update(graph):
    result = graph.invoke({"text": "keep everything"})

    # Each node contributed a partial update; invoke() hands back the merged
    # state, including the input key you passed in.
    assert set(result) == {"text", "cleaned", "word_count", "summary"}
    assert result["text"] == "keep everything"


def test_graph_has_the_expected_nodes(graph):
    nodes = set(graph.get_graph().nodes)

    assert {"clean", "count", "summarize"} <= nodes, (
        f"expected nodes named clean/count/summarize, got {sorted(nodes)}"
    )


def test_graph_is_wired_start_to_end_in_order(graph):
    edges = {(e.source, e.target) for e in graph.get_graph().edges}

    assert (START, "clean") in edges, "START must flow into clean"
    assert ("clean", "count") in edges
    assert ("count", "summarize") in edges
    assert ("summarize", END) in edges, "summarize must flow into END"


def test_nodes_execute_in_dependency_order(graph):
    # stream() yields one chunk per node that ran, in execution order.
    steps = [name for chunk in graph.stream({"text": "a b"}) for name in chunk]

    assert steps == ["clean", "count", "summarize"]
