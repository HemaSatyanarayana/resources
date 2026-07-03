/**
 * MyPromise — a Promises/A+ subset, from scratch.
 *
 * Requirements:
 *   - new MyPromise((resolve, reject) => { ... })  runs the executor immediately;
 *     if it throws, the promise rejects with the thrown value.
 *   - A promise is a state machine: 'pending' -> 'fulfilled' | 'rejected',
 *     and settles at most once.
 *   - .then(onFulfilled, onRejected) returns a NEW MyPromise (chainable). Handlers
 *     run asynchronously (schedule them as microtasks with queueMicrotask).
 *   - If a handler returns a value, the next promise fulfills with it. If it
 *     returns a thenable (including a MyPromise), the chain ADOPTS it (flattens).
 *   - Missing handlers pass the value/reason through the chain.
 *   - .catch(onRejected) === .then(null, onRejected)
 *   - static MyPromise.resolve(v) / MyPromise.reject(r)
 *
 * Because it has a spec-compliant .then, `await myPromise` works.
 */
export class MyPromise {
  constructor(executor) {
    this.state = 'pending'
    this.value = undefined
    this.callbacks = [] // handlers queued while pending

    const resolve = (value) => this._settle('fulfilled', value)
    const reject = (reason) => this._settle('rejected', reason)

    // TODO: run executor(resolve, reject); if it throws, reject with the error.
  }

  _settle(state, value) {
    // TODO:
    //   - if already settled (state !== 'pending'), return.
    //   - THENABLE ADOPTION: if state === 'fulfilled' and `value` is an object or
    //     function with a callable `.then`, call value.then(resolveInner, rejectInner)
    //     to adopt its eventual state (guard against being called twice), and return.
    //   - otherwise set this.state/this.value, then flush queued callbacks as
    //     microtasks and clear the queue.
    throw new Error('TODO: implement _settle')
  }

  then(onFulfilled, onRejected) {
    // TODO: return new MyPromise((resolve, reject) => { ... }) where a `handle`
    // function:
    //   - picks onFulfilled or onRejected based on this.state
    //   - if that handler isn't a function, passes through (resolve or reject
    //     with this.value)
    //   - else tries resolve(handler(this.value)), catching to reject
    // Schedule `handle` as a microtask if already settled, else push to callbacks.
    throw new Error('TODO: implement then')
  }

  catch(onRejected) {
    return this.then(null, onRejected)
  }

  static resolve(value) {
    return new MyPromise((resolve) => resolve(value))
  }

  static reject(reason) {
    return new MyPromise((_, reject) => reject(reason))
  }
}
