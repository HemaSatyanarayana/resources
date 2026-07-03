# 09 — useFetch

Every data-driven round needs this: fetch on mount, model **loading / error / data** as one coherent state, refetch when the input changes, and — the part that separates seniors — **ignore stale responses** so a slow request that resolves after unmount (or after the url changed) doesn't clobber fresh state.

## Concepts

- **Three states, always handled.** `loading`, `error`, and `data` are the minimum. Interviewers deduct points for rendering `data.map(...)` without a loading/error branch.
- **`res.ok` is not automatic.** `fetch` only rejects on network failure; a 404/500 still resolves. You must check `res.ok` and throw yourself.
- **The stale-response guard.** A boolean `ignore`, flipped to `true` in the effect's cleanup, prevents `setState` after unmount and prevents an old url's response from overwriting the new one. (An `AbortController` is the fancier version — mention it.)
- **`[url]` dependency** re-runs the whole thing when the input changes.

## Your task

Implement the effect in [`useFetch.js`](useFetch.js). Tests stub `global.fetch`.

## Run

```bash
npx vitest run exercises/09-use-fetch
```

## Interview follow-ups to expect

- *"Two requests race — the first is slower and resolves last. What renders?"* — without the guard, the stale first response wins. The `ignore` flag (set in cleanup when url changed) discards it. `AbortController.abort()` additionally cancels the in-flight request.
- *"Add retry / caching."* — memoize by url, exponential backoff on error.
- *"Why not fetch in the component body?"* — side effects belong in effects; the body must be pure and can run many times.

## Hints

```js
useEffect(() => {
  let ignore = false
  setState({ data: null, loading: true, error: null })
  fetch(url)
    .then((res) => {
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      return res.json()
    })
    .then((data) => { if (!ignore) setState({ data, loading: false, error: null }) })
    .catch((error) => { if (!ignore) setState({ data: null, loading: false, error }) })
  return () => { ignore = true }
}, [url])
```

<details>
<summary>Reference solution (try first!)</summary>

```jsx
import { useState, useEffect } from 'react'

export function useFetch(url) {
  const [state, setState] = useState({ data: null, loading: true, error: null })

  useEffect(() => {
    let ignore = false
    setState({ data: null, loading: true, error: null })

    fetch(url)
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
        return res.json()
      })
      .then((data) => {
        if (!ignore) setState({ data, loading: false, error: null })
      })
      .catch((error) => {
        if (!ignore) setState({ data: null, loading: false, error })
      })

    return () => {
      ignore = true
    }
  }, [url])

  return state
}
```

</details>
