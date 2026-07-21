"""Make `from graph import ...` mean *this exercise's* graph.py.

Every exercise directory has its own `graph.py`, and every spec imports it as a
top-level module. Python caches modules by name, so a plain `uv run pytest` over
all ten would hand the first `graph` it imported to every spec.

Popping the cache between imports is not enough. These modules use
`from __future__ import annotations`, so a state schema's annotations are
*strings* until something resolves them — which LangGraph does, at
`StateGraph(MyState)` time, by looking up `sys.modules[MyState.__module__]`.
That module name is `"graph"` for all ten exercises. So the mapping has to be
correct while a test *runs*, not just while it is imported.

Two hooks:

  - `pytest_pycollect_makemodule` imports each exercise's `graph.py` under a
    private handle before its spec is collected, so every spec imports its own.
  - `pytest_runtest_setup` re-points `sys.modules["graph"]` at the right one
    before each test, so annotation lookups resolve in the right namespace.

None of this matters when you run one exercise at a time — which is how the
READMEs tell you to work. It exists so `uv run pytest` also does the right
thing.
"""

from __future__ import annotations

import importlib.util
import sys

#: exercise directory -> the module object loaded from its graph.py
_GRAPH_MODULES: dict[str, object] = {}


def _load_graph_module(directory) -> object | None:
    """Import `<directory>/graph.py` as the module named `graph`, once."""
    key = str(directory)
    if key in _GRAPH_MODULES:
        return _GRAPH_MODULES[key]

    graph_py = directory / "graph.py"
    if not graph_py.exists():
        return None

    spec = importlib.util.spec_from_file_location("graph", graph_py)
    module = importlib.util.module_from_spec(spec)

    # Register before executing: a module that imports itself, directly or via a
    # dataclass/TypedDict lookup, must find the same object.
    sys.modules["graph"] = module
    spec.loader.exec_module(module)

    _GRAPH_MODULES[key] = module
    return module


def _activate(directory) -> None:
    key = str(directory)

    while key in sys.path:
        sys.path.remove(key)
    sys.path.insert(0, key)

    module = _GRAPH_MODULES.get(key) or _load_graph_module(directory)
    if module is not None:
        sys.modules["graph"] = module


def pytest_pycollect_makemodule(module_path, parent):
    _activate(module_path.parent)
    return None  # let pytest do the actual collecting


def pytest_runtest_setup(item):
    _activate(item.path.parent)
