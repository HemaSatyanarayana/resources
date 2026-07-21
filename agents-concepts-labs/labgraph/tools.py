"""Toy tools shared by the tool-calling and agent labs.

Deliberately boring and side-effect free so tests can assert on them. The point
of these labs is the *loop* around the tool, not the tool itself.
"""

from __future__ import annotations

from langchain_core.tools import tool

# A tiny stand-in for a search index. Lookup is substring-based and case
# insensitive so a model's phrasing does not have to match exactly.
_FACTS = {
    "langgraph": "LangGraph models agent workflows as a state graph of nodes and edges.",
    "checkpointer": "A checkpointer persists graph state after every step, keyed by thread_id.",
    "reducer": "A reducer decides how a node's update merges into existing state.",
    "interrupt": "interrupt() pauses a graph and surfaces a payload to the caller.",
    "send": "Send() dispatches dynamic parallel work to a node, one task per item.",
}


@tool
def search(query: str) -> str:
    """Look up a fact about LangGraph. Use for questions about how LangGraph works."""
    q = query.lower()
    for key, fact in _FACTS.items():
        if key in q:
            return fact
    return f"No results for {query!r}."


@tool
def add(a: float, b: float) -> float:
    """Add two numbers together."""
    return a + b


@tool
def multiply(a: float, b: float) -> float:
    """Multiply two numbers together."""
    return a * b


@tool
def word_count(text: str) -> int:
    """Count the words in a piece of text."""
    return len(text.split())


#: Convenience bundle for labs that just need "some tools".
ALL_TOOLS = [search, add, multiply, word_count]
