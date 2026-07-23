"""Exercise 03 — Conditional edges: branching and looping.

Straight lines are not agents. An agent decides what to do next, and in
LangGraph that decision is a *conditional edge*: a function that looks at state
and returns the name of the next node.

You are building a draft/review loop — the skeleton of every "critique and
retry" agent.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/03-conditional-routing/README.md

The spec imports these names — they must exist with these exact spellings:

    MAX_REVISIONS       the budget: 3
    DraftState          the state schema
    write_draft         node: first draft
    revise_draft        node: improve the draft, bump the counter
    review              node: record that review happened
    route_after_review  the conditional edge — loop or stop
    build_graph         returns the compiled graph

Then:

    uv run pytest exercises/03-conditional-routing -v
"""
from __future__ import annotations
import operator

from typing import Annotated,TypedDict,Literal
from langgraph.graph import StateGraph, START, END

MAX_REVISIONS = 3

class DraftState(TypedDict):
    topic:	str	
    draft:	str	    
    score:	int	
    revisions:	int	
    history	:Annotated[list[str], operator.add]	

def write_draft(state)-> dict:
    return {
        "draft" : f"draft of {state['topic']}",
        "score" : len(state['topic']) % 11 ,
        "revisions" : 0,
        "history" : ["draft"]
    }

def revise_draft(state)->dict:
    return {
        "draft":f"revision {state['revisions']+1} of {state['topic']}",
        "score": min(10, state['score']+1),
        "revisions": state['revisions'] + 1,
        "history":['revise']
    }

def review(state)->dict:
    return {"history":["review"]}

def route_after_review(state)->Literal['revise','__end__']:
    if state["score"] >= 8 or state['revisions']>=MAX_REVISIONS:
        return END
    return 'revise'

def build_graph():
    graph = StateGraph(DraftState)

    graph.add_node('draft',write_draft)
    graph.add_node('review',review)
    graph.add_node('revise',revise_draft)

    graph.add_edge(START,"draft")
    graph.add_edge('draft','review')
    graph.add_conditional_edges('review',route_after_review)
    graph.add_edge('revise','review')

    return graph.compile()

if __name__ == "__main__":
    from labgraph import print_state

    graph = build_graph()
    for topic in ("ai", "langgraph state machines"):
        result = graph.invoke(
            {"topic": topic, "draft": "", "score": 0, "revisions": 0, "history": []}
        )
        print_state(result, title=f"topic={topic!r}")