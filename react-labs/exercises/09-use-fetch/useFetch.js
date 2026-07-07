import { useState, useEffect } from "react";

/**
 * useFetch — data fetching with loading / error / data states and cleanup.
 *
 * Returns { data, loading, error }.
 *
 * Behavior:
 *   - On mount (and whenever `url` changes): set loading=true, clear old data/error,
 *     then fetch(url).
 *   - If the response is not ok (res.ok === false), treat it as an error:
 *     throw new Error(`HTTP ${res.status}`).
 *   - On success: set { data: <parsed json>, loading: false, error: null }.
 *   - On failure: set { data: null, loading: false, error: <the Error> }.
 *   - CANCELLATION: if the component unmounts (or url changes) before the
 *     request resolves, do NOT call setState — a stale response must be ignored.
 *     Use an `ignore` flag captured by the effect's cleanup.
 */
export function useFetch(url) {
  const [state, setState] = useState({
    data: null,
    loading: true,
    error: null,
  });

  useEffect(() => {
    // TODO:
    //   let ignore = false
    //   setState(loading)
    //   fetch(url) -> check res.ok -> res.json() -> setState(data) unless ignore
    //   .catch -> setState(error) unless ignore
    //   return () => { ignore = true }
    let ignore = false;
    setState({ data: null, loading: true, error: null });

    fetch(url)
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        if (!ignore) {
          setState({
            data,
            loading: false,
            error: null,
          });
        }
      })
      .catch((err) => {
        if (!ignore) {
          setState({
            data: null,
            loading: false,
            error: err,
          });
        }
      });

    return () => {
      ignore = true;
    };
  }, [url]);

  return state;
}
