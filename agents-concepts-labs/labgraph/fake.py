"""A deterministic chat model for tests.

Real models are non-deterministic, cost money, and need network. Neither is
useful when the thing under test is *your graph wiring*. `ScriptedChatModel`
hands back a fixed list of replies in order, so a test can assert exactly what
your graph did with them.

It behaves like a real chat model where it matters:

* `bind_tools(...)` works (and records what was bound), so agent code written
  against a real provider runs unchanged.
* replies may carry `tool_calls`, so you can drive a full tool-calling loop.
* `_stream` emits word-by-word chunks, so `stream_mode="messages"` is meaningful.

Typical use::

    model = ScriptedChatModel(["hello", tool_call("search", {"q": "cats"}), "done"])
    ...
    assert model.call_count == 3
    assert "cats" in model.last_prompt_text()
"""

from __future__ import annotations

import re
from typing import Any, Iterator, Sequence

from langchain_core.callbacks import CallbackManagerForLLMRun
from langchain_core.language_models import BaseChatModel
from langchain_core.messages import AIMessage, AIMessageChunk, BaseMessage
from langchain_core.outputs import ChatGeneration, ChatGenerationChunk, ChatResult
from pydantic import Field, PrivateAttr


def tool_call(name: str, args: dict[str, Any], *, id: str | None = None) -> AIMessage:
    """An assistant reply that asks for one tool call.

    Use it inside a `ScriptedChatModel` script to make the model "decide" to call
    a tool::

        ScriptedChatModel([tool_call("add", {"a": 2, "b": 3}), "The answer is 5"])
    """
    return AIMessage(
        content="",
        tool_calls=[{"name": name, "args": args, "id": id or f"call_{name}"}],
    )


def tool_calls(*calls: tuple[str, dict[str, Any]]) -> AIMessage:
    """An assistant reply asking for several tool calls at once (parallel calls)."""
    return AIMessage(
        content="",
        tool_calls=[
            {"name": name, "args": args, "id": f"call_{name}_{i}"}
            for i, (name, args) in enumerate(calls)
        ],
    )


class ScriptedChatModel(BaseChatModel):
    """Replays `responses` in order, one per invocation.

    Args:
        responses: strings (become plain `AIMessage`s) or ready-made `AIMessage`s.
        loop: when True, restart the script instead of raising once exhausted.
              Handy for agent loops whose length you do not want to predict.
    """

    responses: list[Any] = Field(default_factory=list)
    loop: bool = False

    _index: int = PrivateAttr(default=0)
    _calls: list[list[BaseMessage]] = PrivateAttr(default_factory=list)
    _bound_tools: list[Any] = PrivateAttr(default_factory=list)

    def __init__(self, responses: Sequence[Any] | None = None, **kwargs: Any) -> None:
        # Allow the natural `ScriptedChatModel(["a", "b"])` in addition to the
        # keyword form pydantic would otherwise force on us.
        super().__init__(responses=list(responses or []), **kwargs)

    @property
    def _llm_type(self) -> str:
        return "scripted-chat-model"

    # --- inspection helpers used by tests ---------------------------------

    @property
    def call_count(self) -> int:
        """How many times the model was invoked."""
        return len(self._calls)

    @property
    def calls(self) -> list[list[BaseMessage]]:
        """The message list passed on each invocation, in order."""
        return self._calls

    @property
    def bound_tools(self) -> list[Any]:
        """Tools handed to `bind_tools`, if it was called."""
        return self._bound_tools

    def last_prompt(self) -> list[BaseMessage]:
        """Messages from the most recent invocation."""
        if not self._calls:
            raise AssertionError("model was never called")
        return self._calls[-1]

    def last_prompt_text(self) -> str:
        """Most recent invocation's messages flattened to one lowercase string."""
        return "\n".join(str(m.content) for m in self.last_prompt()).lower()

    def reset(self) -> None:
        """Rewind the script and forget recorded calls."""
        self._index = 0
        self._calls = []

    # --- chat model interface ---------------------------------------------

    def bind_tools(self, tools: Sequence[Any], **kwargs: Any) -> "ScriptedChatModel":
        """Record the tools and return self.

        A real provider would serialise these into the request. The script
        already decides what gets called, so recording is enough — but your graph
        code can call `bind_tools` exactly as it would in production.
        """
        self._bound_tools = list(tools)
        return self

    def _next_message(self, messages: list[BaseMessage]) -> AIMessage:
        self._calls.append(list(messages))

        if self._index >= len(self.responses):
            if self.loop and self.responses:
                self._index = 0
            else:
                raise AssertionError(
                    f"ScriptedChatModel ran out of responses: the graph called the "
                    f"model {self.call_count} time(s) but the script only has "
                    f"{len(self.responses)}. Either the graph looped more than "
                    f"expected, or the script is too short."
                )

        reply = self.responses[self._index]
        self._index += 1
        return AIMessage(content=reply) if isinstance(reply, str) else reply

    def _generate(
        self,
        messages: list[BaseMessage],
        stop: list[str] | None = None,
        run_manager: CallbackManagerForLLMRun | None = None,
        **kwargs: Any,
    ) -> ChatResult:
        return ChatResult(generations=[ChatGeneration(message=self._next_message(messages))])

    def _stream(
        self,
        messages: list[BaseMessage],
        stop: list[str] | None = None,
        run_manager: CallbackManagerForLLMRun | None = None,
        **kwargs: Any,
    ) -> Iterator[ChatGenerationChunk]:
        message = self._next_message(messages)

        # Tool calls have no text to split — emit the message as a single chunk.
        if message.tool_calls or not isinstance(message.content, str):
            yield ChatGenerationChunk(
                message=AIMessageChunk(
                    content=message.content, tool_calls=message.tool_calls
                )
            )
            return

        # Split on whitespace but keep it, so rejoining the chunks reproduces the
        # original text exactly — the property a streaming test wants to assert.
        for piece in re.findall(r"\S+\s*", message.content) or [""]:
            chunk = ChatGenerationChunk(message=AIMessageChunk(content=piece))
            if run_manager:
                run_manager.on_llm_new_token(piece, chunk=chunk)
            yield chunk
