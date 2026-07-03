import { useState, useRef, useEffect } from 'react'

/**
 * Stopwatch — useRef for the interval id, no stale closures, effect cleanup.
 *
 * Required UI (the tests rely on this):
 *   - an element with data-testid="elapsed" showing seconds to one decimal,
 *     e.g. "0.0", "1.0" (compute as (ms / 1000).toFixed(1))
 *   - a button "Start"  (disabled while running)
 *   - a button "Stop"   (disabled while stopped)
 *   - a button "Reset"
 *
 * Rules:
 *   - Start begins ticking every 100ms, adding 100 to the elapsed ms.
 *   - Stop pauses; the elapsed value freezes.
 *   - Reset stops and returns elapsed to 0.
 *   - The interval MUST be cleared on Stop/Reset and on unmount (no leaks).
 *   - Use a functional update inside the tick so you never read stale `ms`.
 */
export default function Stopwatch() {
  const [ms, setMs] = useState(0)
  const [running, setRunning] = useState(false)
  const intervalRef = useRef(null)

  const start = () => {
    // TODO: if already running, do nothing. Otherwise set running true and
    // start an interval (100ms) that does setMs(m => m + 100). Save its id in
    // intervalRef.current.
  }

  const stop = () => {
    // TODO: set running false, clearInterval(intervalRef.current), null it out.
  }

  const reset = () => {
    // TODO: stop, then setMs(0).
  }

  useEffect(() => {
    // TODO: cleanup on unmount — clear any live interval.
  }, [])

  return (
    <div>
      <output data-testid="elapsed">{(ms / 1000).toFixed(1)}</output>
      <button onClick={start} disabled={running}>
        Start
      </button>
      <button onClick={stop} disabled={!running}>
        Stop
      </button>
      <button onClick={reset}>Reset</button>
    </div>
  )
}
