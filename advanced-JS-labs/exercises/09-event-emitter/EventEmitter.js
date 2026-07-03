/**
 * EventEmitter — a minimal pub/sub.
 *
 *   on(event, cb)   subscribe; RETURN an unsubscribe function
 *   off(event, cb)  remove a specific listener
 *   once(event, cb) subscribe, but auto-remove after the first emit
 *   emit(event, ...args)  call every listener for `event` with args
 *
 * A listener registered twice, then removed once, should still fire (or use a
 * Set so duplicates collapse — your call, but be consistent). These tests use
 * distinct callbacks, so a Set-per-event works cleanly.
 */
export class EventEmitter {
  constructor() {
    this.listeners = new Map() // event -> Set<cb>
  }

  on(event, cb) {
    // TODO: ensure a Set exists for `event`, add cb, return () => this.off(event, cb)
    throw new Error('TODO: implement on')
  }

  off(event, cb) {
    // TODO: remove cb from the event's Set (if any)
    throw new Error('TODO: implement off')
  }

  once(event, cb) {
    // TODO: register a wrapper that removes itself, then calls cb(...args).
    // Return the unsubscribe handle too.
    throw new Error('TODO: implement once')
  }

  emit(event, ...args) {
    // TODO: call each listener for `event` with args. Iterate over a copy so a
    // listener that unsubscribes during emit doesn't corrupt the iteration.
    throw new Error('TODO: implement emit')
  }
}
