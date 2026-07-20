// Package pubsub — Lab 03: an in-memory, generic publish/subscribe broker.
//
// Read README.md first; broker_test.go is the spec. Fill in every TODO.
// Run: go test -race -v ./03-pubsub/
package pubsub

// Broker routes messages of type T from publishers to per-topic subscribers.
type Broker[T any] struct {
	// TODO: fields. You'll need:
	//   - a sync.RWMutex (Publish is read-mostly; Subscribe/cancel/Close write),
	//   - a map from topic -> set of subscriber channels (map[string]map[int]chan T
	//     works well, using an int id so cancel can find & remove its own channel),
	//   - a counter for the next subscriber id,
	//   - a `closed bool` flag.
}

// NewBroker creates an empty broker.
func NewBroker[T any]() *Broker[T] {
	// TODO: initialise the maps.
	panic("TODO: implement NewBroker")
}

// Subscribe registers interest in topic and returns a receive-only channel of
// messages plus a cancel function. Calling cancel unsubscribes and closes the
// channel; it is safe to call cancel more than once. If the broker is already
// closed, Subscribe returns an already-closed channel and a no-op cancel.
func (b *Broker[T]) Subscribe(topic string) (<-chan T, func()) {
	// TODO:
	//   - Lock. If closed: return a closed channel + no-op cancel.
	//   - Allocate an id and a BUFFERED channel (e.g. cap 16) so a briefly-slow
	//     subscriber doesn't drop everything.
	//   - Store it under subs[topic][id].
	//   - Build a cancel closure that locks, deletes subs[topic][id] if present,
	//     and closes that channel (guard so a second call is a no-op).
	panic("TODO: implement Subscribe")
}

// Publish delivers msg to every current subscriber of topic. Delivery is
// non-blocking per subscriber: if a subscriber's buffer is full, msg is dropped
// for that subscriber (a slow consumer must never block the publisher or others).
// Publishing to a closed broker, or a topic with no subscribers, is a no-op.
func (b *Broker[T]) Publish(topic string, msg T) {
	// TODO:
	//   - RLock. If closed: return.
	//   - For each subscriber channel of `topic`: non-blocking send
	//     (select { case ch <- msg: default: }).
	panic("TODO: implement Publish")
}

// Close shuts the broker down: it closes every subscriber channel and rejects
// future publishes. It is safe to call more than once.
func (b *Broker[T]) Close() {
	// TODO:
	//   - Lock. If already closed: return. Set closed = true.
	//   - Close every subscriber channel and clear the maps.
	panic("TODO: implement Close")
}
