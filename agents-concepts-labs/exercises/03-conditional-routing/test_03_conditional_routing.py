"""Spec for exercise 03. Do not edit — make it pass."""

from __future__ import annotations

import pytest
from langgraph.errors import GraphRecursionError
from langgraph.graph import END

from graph import MAX_REVISIONS, build_graph, review, revise_draft, route_after_review, write_draft


def initial(topic: str) -> dict:
    return {"topic": topic, "draft": "", "score": 0, "revisions": 0, "history": []}


# --- nodes ----------------------------------------------------------------


def test_write_draft_scores_deterministically():
    update = write_draft(initial("langgraph"))

    assert update["draft"] == "draft of langgraph"
    assert update["score"] == 9  # len("langgraph") % 11
    assert update["revisions"] == 0
    assert update["history"] == ["draft"]


def test_revise_bumps_score_and_revision_count():
    update = revise_draft({"topic": "ai", "draft": "d", "score": 3, "revisions": 1, "history": []})

    assert update["score"] == 4
    assert update["revisions"] == 2
    assert update["draft"] == "revision 2 of ai"


def test_revise_caps_the_score_at_ten():
    update = revise_draft({"topic": "ai", "draft": "d", "score": 10, "revisions": 0, "history": []})

    assert update["score"] == 10


def test_review_only_records_history():
    assert review(initial("x")) == {"history": ["review"]}


# --- the router -----------------------------------------------------------


def test_router_stops_when_quality_is_good_enough():
    assert route_after_review({"score": 8, "revisions": 0}) == END


def test_router_loops_when_quality_is_low_and_budget_remains():
    assert route_after_review({"score": 2, "revisions": 0}) == "revise"


def test_router_stops_when_the_revision_budget_is_spent():
    # Even though the draft is still bad, we are out of attempts.
    assert route_after_review({"score": 2, "revisions": MAX_REVISIONS}) == END


# --- the loop end to end --------------------------------------------------


@pytest.fixture
def graph():
    return build_graph()


def test_good_first_draft_skips_revision(graph):
    # "langgraph" scores 9 -> already above the bar.
    result = graph.invoke(initial("langgraph"))

    assert result["revisions"] == 0
    assert result["history"] == ["draft", "review"]


def test_mediocre_draft_is_revised_until_it_passes(graph):
    # "prompts" scores 7 -> one revision takes it to 8 and it exits.
    result = graph.invoke(initial("prompts"))

    assert result["score"] == 8
    assert result["revisions"] == 1
    assert result["history"] == ["draft", "review", "revise", "review"]


def test_hopeless_draft_stops_at_the_revision_budget(graph):
    # "ai" scores 2 and gains 1 per revision, so it can never reach 8.
    # The budget — not the quality bar — is what ends this run.
    result = graph.invoke(initial("ai"))

    assert result["revisions"] == MAX_REVISIONS
    assert result["score"] < 8, "this run should end unhappy, not by passing review"
    assert result["history"].count("revise") == MAX_REVISIONS


def test_graph_reports_a_conditional_edge_out_of_review(graph):
    drawn = graph.get_graph()
    conditional = {(e.source, e.target) for e in drawn.edges if e.conditional}

    assert ("review", "revise") in conditional, (
        "review -> revise must come from add_conditional_edges, not a plain edge"
    )


def test_recursion_limit_is_the_backstop_against_runaway_loops(graph):
    # Your MAX_REVISIONS budget is the *intended* stop. recursion_limit is the
    # framework's hard ceiling on supersteps, and it raises rather than hanging.
    # Squeeze it below what the "ai" run needs and LangGraph aborts.
    with pytest.raises(GraphRecursionError):
        graph.invoke(initial("ai"), {"recursion_limit": 4})
