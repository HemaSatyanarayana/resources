# React Labs — Master the Machine Coding Round

A hands-on, **test-driven** curriculum for the React machine coding interview. These are the components interviewers actually ask you to build in 30–45 minutes. Each exercise ships with:

- **A README** explaining the concept, the exact task, the interview follow-ups you'll get, and what "good React" looks like.
- **Starter code** with the component signature and `TODO` bodies for you to implement.
- **A test suite** written with React Testing Library — implement until `npm test` is green.
- **A reference solution** hidden in a `<details>` block in each README (peek only after trying).

The whole repo is a single package. Every exercise lives in its own folder, so you can work on them independently.

## Setup

```bash
cd react-labs
npm install
```

## How to use this

1. Pick an exercise directory, e.g. `exercises/01-counter`.
2. Read its `README.md`.
3. Open the starter file (e.g. `Counter.jsx`) and replace each `TODO`.
4. Run the tests for that exercise:

   ```bash
   npx vitest run exercises/01-counter
   ```

5. Green? Move on. Red? Read the failure, fix, repeat.

Run **everything** at once:

```bash
npm test
```

Watch mode (re-runs on save — the way to actually work through these):

```bash
npm run test:watch
```

## Curriculum

### Track A — React core (01–05)

The state-and-events fundamentals every round is built on. Controlled inputs, list rendering, derived state, lifting state up.

| #  | Exercise | What you'll master |
|----|----------|--------------------|
| 01 | [Counter](exercises/01-counter) | `useState`, functional updates, event handlers, clamping |
| 02 | [Todo List](exercises/02-todo-list) | Controlled inputs, immutable list updates, `key`, filtering |
| 03 | [Star Rating](exercises/03-star-rating) | Props, hover vs selected state, derived rendering, reuse |
| 04 | [Accordion](exercises/04-accordion) | Conditional rendering, single vs multi-open, ARIA |
| 05 | [Tabs](exercises/05-tabs) | Active index state, keyboard nav, `role="tab"` semantics |

### Track B — Hooks mastery (06–08)

Custom hooks, timers, refs, and cleanup — the topics that separate "can use React" from "understands React."

| #  | Exercise | What you'll master |
|----|----------|--------------------|
| 06 | [Custom Hooks](exercises/06-custom-hooks) | `useToggle`, `usePrevious`, `useLocalStorage` — composition & rules of hooks |
| 07 | [useDebounce](exercises/07-use-debounce) | `setTimeout`, effect cleanup, stale-closure traps |
| 08 | [Stopwatch](exercises/08-stopwatch) | `useRef` for interval ids, start/stop/reset, no stale state |

### Track C — Async & data (09–12)

Fetching, loading/error states, and the async UI patterns interviewers love: debounced search, pagination, infinite scroll.

| #  | Exercise | What you'll master |
|----|----------|--------------------|
| 09 | [useFetch](exercises/09-use-fetch) | Data-fetching hook, loading/error/data, cancellation on unmount |
| 10 | [Autocomplete](exercises/10-autocomplete) | Debounced search, keyboard navigation, race conditions |
| 11 | [Pagination](exercises/11-pagination) | Page windowing, boundary logic, page-size math |
| 12 | [Infinite Scroll](exercises/12-infinite-scroll) | `IntersectionObserver`, appending pages, sentinel elements |

### Track D — Interaction & advanced (13–15)

Portals, recursion, and game logic — the "senior" flavored questions.

| #  | Exercise | What you'll master |
|----|----------|--------------------|
| 13 | [Modal](exercises/13-modal) | `createPortal`, Escape-to-close, backdrop clicks, focus |
| 14 | [Nested Checkboxes](exercises/14-nested-checkboxes) | Recursive components, parent/child indeterminate state |
| 15 | [Tic Tac Toe](exercises/15-tic-tac-toe) | Derived winner logic, immutable board, turn state |

## Recommended order

Go in numeric order — later exercises assume earlier concepts. Track A is non-negotiable warm-up; do not skip it even if it feels easy, because the interview *starts* here and you want the muscle memory. Track B's cleanup discipline is a prerequisite for the async work in Track C.

## The machine-coding rubric (how you're actually judged)

Interviewers score more than "does it work." Practice hitting all of these:

- **Correctness** — the happy path works, and the obvious edge cases (empty, one item, boundaries) work too.
- **State modeling** — minimal state; everything else is *derived*. If you can compute it, don't store it.
- **Immutability** — never mutate state. New array/object every update (`[...items]`, `{...obj}`).
- **Controlled components** — inputs driven by state, not the DOM.
- **Keys** — stable, unique `key` on lists (never the array index when items reorder/delete).
- **Cleanup** — every `setTimeout`/`setInterval`/subscription/`addEventListener` is torn down in the effect's cleanup.
- **Accessibility** — real semantics: `<button>` not `<div onClick>`, labels on inputs, ARIA roles on widgets.
- **Talking while coding** — narrate your state shape and trade-offs. Silence loses points.

## Idioms this lab drills into you

- **Derive, don't store.** The winner of a tic-tac-toe game is a function of the board, not a piece of state.
- **Updater functions for state that depends on previous state:** `setCount(c => c + 1)`, never `setCount(count + 1)` in a loop/async path.
- **Clean up every effect.** A `setInterval` without a `clearInterval` is a bug and a red flag.
- **Lift state up** to the closest common parent; pass data down as props, changes up as callbacks.
- **Accessibility is not optional.** Semantic HTML is faster to write *and* scores better.

## Toolbelt

```bash
npm test                                  # run all exercises once
npm run test:watch                        # re-run on save
npx vitest run exercises/03-star-rating   # one exercise
npx vitest run -t "clamps"                # tests matching a name
```

Happy hacking. ⚛️
