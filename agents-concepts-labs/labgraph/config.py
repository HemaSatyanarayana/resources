"""Provider-agnostic model configuration.

Every lab gets its chat model from here, and this module reads *only* environment
variables. No provider is hardcoded in the course, so the same graphs run against
OpenRouter, Gemini's OpenAI-compatible endpoint, a local Ollama, or anything else
that speaks the OpenAI wire format.

    LLM_BASE_URL   e.g. https://openrouter.ai/api/v1
    LLM_API_KEY    provider key
    LLM_MODEL      e.g. anthropic/claude-sonnet-4.5
    LLM_TEMPERATURE / LLM_MAX_TOKENS / LLM_TIMEOUT   optional

Tests never call this — they use `labgraph.ScriptedChatModel` so the suite is
deterministic and runs with no key and no network.
"""

from __future__ import annotations

import os
from typing import Any

from dotenv import load_dotenv
from langchain_core.language_models import BaseChatModel

# Load `.env` from the repo root if present. override=False so a real exported
# environment variable always wins over the file.
load_dotenv(override=False)

REQUIRED_VARS = ("LLM_BASE_URL", "LLM_API_KEY", "LLM_MODEL")

_SETUP_HINT = """\
No LLM is configured. Copy `.env.example` to `.env` and fill in:

    LLM_BASE_URL=https://openrouter.ai/api/v1
    LLM_API_KEY=sk-or-v1-...
    LLM_MODEL=anthropic/claude-sonnet-4.5

Missing: {missing}

You do not need this to work through the exercises — `pytest` runs entirely on a
fake model. This is only for running a graph against a real provider.\
"""


def missing_vars() -> list[str]:
    """Which required env vars are unset or empty."""
    return [v for v in REQUIRED_VARS if not os.environ.get(v)]


def llm_configured() -> bool:
    """True when a real model can be built. Use this to skip live tests."""
    return not missing_vars()


def _float_env(name: str) -> float | None:
    raw = os.environ.get(name)
    return float(raw) if raw else None


def _int_env(name: str) -> int | None:
    raw = os.environ.get(name)
    return int(raw) if raw else None


def get_chat_model(**overrides: Any) -> BaseChatModel:
    """Build a chat model from the environment.

    Any keyword argument overrides the env-derived value, e.g.
    `get_chat_model(temperature=0.7)`.

    Raises RuntimeError with setup instructions when the env is incomplete.
    """
    missing = missing_vars()
    if missing:
        raise RuntimeError(_SETUP_HINT.format(missing=", ".join(missing)))

    # Imported lazily so the rest of the course works even if the OpenAI client
    # extra is not installed.
    from langchain_openai import ChatOpenAI

    kwargs: dict[str, Any] = {
        "base_url": os.environ["LLM_BASE_URL"],
        "api_key": os.environ["LLM_API_KEY"],
        "model": os.environ["LLM_MODEL"],
        "temperature": _float_env("LLM_TEMPERATURE") or 0.0,
    }
    if (max_tokens := _int_env("LLM_MAX_TOKENS")) is not None:
        kwargs["max_tokens"] = max_tokens
    if (timeout := _float_env("LLM_TIMEOUT")) is not None:
        kwargs["timeout"] = timeout

    kwargs.update(overrides)
    return ChatOpenAI(**kwargs)
