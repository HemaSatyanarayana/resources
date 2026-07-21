"""The marker every unfinished scaffold raises."""

from __future__ import annotations


class TODO(NotImplementedError):
    """Raised by scaffold code you have not implemented yet.

    Exercise files import this and raise it from every stub, so the module still
    imports cleanly (and the whole suite still collects) while the tests fail
    loudly and specifically until you write the real thing.
    """

    def __init__(self, hint: str = "not implemented yet") -> None:
        super().__init__(f"TODO: {hint}")
