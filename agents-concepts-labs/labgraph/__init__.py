"""Shared helpers for the LangGraph labs.

This package is *already written* — it is scaffolding, not an exercise. You will
import from it constantly:

    from labgraph import get_chat_model          # env-configured real model
    from labgraph import ScriptedChatModel, tool_call   # deterministic fake
    from labgraph import search, add, multiply          # toy tools
"""

from labgraph.config import get_chat_model, llm_configured, missing_vars
from labgraph.fake import ScriptedChatModel, tool_call, tool_calls
from labgraph.pretty import print_messages, print_state
from labgraph.todo import TODO
from labgraph.tools import ALL_TOOLS, add, multiply, search, word_count

__all__ = [
    "TODO",
    "get_chat_model",
    "llm_configured",
    "missing_vars",
    "ScriptedChatModel",
    "tool_call",
    "tool_calls",
    "print_messages",
    "print_state",
    "ALL_TOOLS",
    "add",
    "multiply",
    "search",
    "word_count",
]
