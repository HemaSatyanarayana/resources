"""Printing helpers, for when you run a lab with `python graph.py` and want to
see what actually happened."""

from __future__ import annotations

from typing import Any, Iterable, Mapping

from langchain_core.messages import BaseMessage

_LABELS = {"human": "you", "ai": "model", "tool": "tool", "system": "system"}


def print_messages(messages: Iterable[BaseMessage]) -> None:
    """Render a message list compactly, including tool calls."""
    for m in messages:
        label = _LABELS.get(m.type, m.type)
        content = str(m.content).strip()
        if content:
            print(f"  [{label}] {content}")
        for call in getattr(m, "tool_calls", []) or []:
            print(f"  [{label}] -> {call['name']}({call['args']})")


def print_state(state: Mapping[str, Any], *, title: str = "state") -> None:
    """Render a state dict, expanding a `messages` key if there is one."""
    print(f"\n== {title} ==")
    for key, value in state.items():
        if key == "messages" and isinstance(value, list):
            print("  messages:")
            print_messages(value)
        else:
            print(f"  {key}: {value!r}")
