/**
 * promisify(fn)
 * Convert a Node-style callback function — one whose LAST argument is a
 * callback (err, result) => ... — into one that returns a Promise.
 *
 *   const readFileP = promisify(fs.readFile)
 *   await readFileP('a.txt', 'utf8')
 */
export function promisify(fn) {
  // TODO: return (...args) => new Promise((resolve, reject) => {
  //   call fn with args + a (err, result) callback.
  // })
  throw new Error('TODO: implement promisify')
}

/**
 * delay(ms, value)
 * Resolve with `value` after `ms` milliseconds.
 */
export function delay(ms, value) {
  // TODO
  throw new Error('TODO: implement delay')
}

/**
 * series(tasks)
 * `tasks` is an array of thunks that each return a Promise. Run them ONE AT A
 * TIME, in order, and resolve with the array of their results (in order).
 * Contrast with Promise.all, which runs them concurrently.
 */
export function series(tasks) {
  // TODO: chain with reduce, starting from Promise.resolve([]).
  throw new Error('TODO: implement series')
}
