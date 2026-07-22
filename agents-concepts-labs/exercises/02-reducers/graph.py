"""Exercise 02 — Reducers: controlling how updates merge.

In exercise 01 every node owned its own key, so "the update overwrites the old
value" was fine. The moment two nodes write the *same* key — or one node writes
it repeatedly in a loop — you have to say what merging means. That is a reducer.

WRITE THIS FILE YOURSELF, from the imports down. It is empty on purpose.

The README explains every import you need, the concepts behind them, and exactly
what to build:

    exercises/02-reducers/README.md

The spec imports these names — they must exist with these exact spellings:

    keep_max       a custom reducer: (current, update) -> larger
    RunState       state schema with four different merge behaviours
    stage_one      node writing all four keys
    stage_two      node writing all four keys again
    build_graph    returns the compiled graph

Then:

    uv run pytest exercises/02-reducers -v
"""

from __future__ import annotations
import operator

from typing import Annotated,TypedDict

from langchain_core.messages import HumanMessage,AIMessage
from langgraph.graph import START,END,StateGraph,add_messages

def keep_max(current: int | None, update: int)-> int:
    if  current is None:
        return update
    return max(current,update)

class RunState(TypedDict):
    status: str
    log: Annotated[list[str],operator.add]
    messages:Annotated[list,add_messages]
    high_score:Annotated[int,keep_max]


def stage_one(state:RunState)->dict:
    return {
        "status"	:"one",
        "log"	:["entered one"],
        "messages"	:[AIMessage("hello from one")],
        "high_score"	:10
    }

def stage_two(state:RunState)->dict:
    return {
        "status"	:"two",
        "log"	:["entered two"],
        "messages"	:[AIMessage("hello from two")],
        "high_score"	:5
    }


def build_graph():
    graph = StateGraph(RunState)

    graph.add_node("one",stage_one)
    graph.add_node("two",stage_two)

    graph.add_edge(START,"one")
    graph.add_edge("one","two")
    graph.add_edge("two",END)

    return graph.compile()

if __name__ == "__main__":
    from labgraph import print_state

    graph = build_graph()
    result = graph.invoke(
        {
            "status": "start",
            "log": ["initial"],
            "messages": [HumanMessage("hi")],
            "high_score": 0,
        }
    )
    print_state(result, title="after two stages")
    print("\nNotice: status was overwritten, log/messages accumulated,")
    print("and high_score kept 10 even though stage_two wrote 5.")
