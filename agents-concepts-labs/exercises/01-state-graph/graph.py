"""Exercise 01 — Your first StateGraph.

Build a three-node text pipeline. No LLM involved: the point is to internalise
the execution model before anything non-deterministic enters the picture.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/01-state-graph/README.md

The spec imports these names — they must exist with these exact spellings:

    PipelineState    the state schema
    clean_text       node: normalise `text` into `cleaned`
    count_words      node: count words of `cleaned` into `word_count`
    summarize        node: format `summary`
    build_pipeline   returns the compiled graph

Then:

    uv run pytest exercises/01-state-graph -v
"""

from __future__ import annotations
from typing import TypedDict
from langgraph.graph import StateGraph,START,END


class PipelineState(TypedDict):
    text:str
    cleaned:str
    word_count:int
    summary:str


def clean_text(state)->dict:
    lower_text = " ".join(state["text"].lower().split())
    return {"cleaned": lower_text}

def count_words(state)->dict:
    return {"word_count":len(state["cleaned"].split())}

def summarize(state)->dict:
    return {"summary": f"{state['word_count']} words: {state['cleaned']}"}

def build_pipeline():
    builder =   StateGraph(PipelineState)

    builder.add_node("clean",clean_text)
    builder.add_node("count",count_words)
    builder.add_node("summarize",summarize)

    builder.add_edge(START,"clean")
    builder.add_edge("clean","count")
    builder.add_edge("count","summarize")
    builder.add_edge("summarize",END)

    return builder.compile()

if __name__ == "__main__":
    from labgraph import print_state

    graph = build_pipeline()
    result = graph.invoke({"text": "  LangGraph makes agent  control flow EXPLICIT  "})
    print_state(result, title="final state")