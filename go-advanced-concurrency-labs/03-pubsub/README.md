# 03 — In-Memory Pub/Sub Broker

A publish/subscribe broker decouples producers from consumers: publishers send
to a **topic**, and every subscriber of that topic gets a copy. It's the core of
event buses, websocket hubs, and in-process message fan-out. You'll build a
**generic** one (`Broker[T]`) and, along the way, confront the single trickiest
bug in channel code: **closing a channel someone else might be sending on.**

## The system

```
                       topic "news"
 Publish("news", m) ──▶ ┌─ subscriber A  (buffered chan, cap 16)
                        ├─ subscriber B
                        └─ subscriber C
                       topic "sports"
 Publish("sports", n)──▶ └─ subscriber D
```

Operations:

- **`Subscribe(topic)`** → `(<-chan T, cancel)`. Registers a fresh buffered
  channel and returns it plus a `cancel` func that unsubscribes and closes it.
- **`Publish(topic, msg)`** → fans `msg` out to every current subscriber of that
  topic, **non-blocking**: a full subscriber buffer means the message is *dropped*
  for that subscriber. A slow consumer must never stall the publisher or its peers.
- **`Close()`** → closes every subscriber channel and rejects further publishes.

### The concurrency design

Use one `sync.RWMutex` guarding a `map[string]map[int]chan T` (topic → id →
channel). `Publish` is the hot path and only *reads* the structure, so it takes
`RLock`; `Subscribe`, `cancel`, and `Close` mutate it under the full `Lock`.

**Why the int id?** `cancel` needs to remove *its own* channel from the topic's
set and close it. Channels aren't comparable-friendly as map keys in practice
(and you can't easily find "which one is mine"), so key subscribers by a
monotonic id the broker hands out.

### The close-safety rule (internalize this)

> Sending on a closed channel **panics**. A channel must be closed by exactly one
> party, and never while a send might be in flight.

Here's how the mutex makes that safe:

- `Publish` holds `RLock` while sending. `cancel` and `Close` hold the full
  `Lock` to close channels. `RLock` and `Lock` are mutually exclusive, so **no
  close can happen while any send is in progress.** That single invariant is
  what prevents the panic.
- `cancel` must be idempotent: on the second call the id is already gone from the
  map, so it closes nothing. Don't close a channel you didn't just delete.

## Your task

Implement, from scratch, in `broker.go` (package `pubsub`):

```go
type Broker[T any] struct { /* your fields */ }

func NewBroker[T any]() *Broker[T]
func (b *Broker[T]) Subscribe(topic string) (<-chan T, func())
func (b *Broker[T]) Publish(topic string, msg T)
func (b *Broker[T]) Close()
```

Subscriber channels are **buffered** (cap 16 is fine). Delivery is drop-on-full.
`cancel` and `Close` are both idempotent. Subscribing to a closed broker returns
an already-closed channel and a no-op cancel.

## Run

```bash
go test -race -v ./03-pubsub/
```

## Hints

- Non-blocking send: `select { case ch <- msg: default: }`.
- Give each subscriber an id from a counter you bump under the lock.
- `cancel` closure captures `topic` and `id`; it locks, checks the channel is
  still in the map, deletes it, then closes it.
- `Close`: set `closed = true`, then range every topic's map closing each channel.
- Publishing/subscribing after `Close` must be a safe no-op, not a panic.

<details>
<summary>Reference solution</summary>

```go
package pubsub

import "sync"

type Broker[T any] struct {
	mu     sync.RWMutex
	subs   map[string]map[int]chan T
	nextID int
	closed bool
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{subs: make(map[string]map[int]chan T)}
}

func (b *Broker[T]) Subscribe(topic string) (<-chan T, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		ch := make(chan T)
		close(ch)
		return ch, func() {}
	}

	id := b.nextID
	b.nextID++
	ch := make(chan T, 16)
	if b.subs[topic] == nil {
		b.subs[topic] = make(map[int]chan T)
	}
	b.subs[topic][id] = ch

	cancel := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if m, ok := b.subs[topic]; ok {
			if c, ok := m[id]; ok {
				delete(m, id)
				close(c)
			}
		}
	}
	return ch, cancel
}

func (b *Broker[T]) Publish(topic string, msg T) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return
	}
	for _, ch := range b.subs[topic] {
		select {
		case ch <- msg:
		default: // subscriber buffer full — drop
		}
	}
}

func (b *Broker[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for _, m := range b.subs {
		for id, ch := range m {
			delete(m, id)
			close(ch)
		}
	}
}
```

The whole design rests on one line of reasoning: closes happen under `Lock`,
sends happen under `RLock`, and those never overlap — so no send ever races a
close.

</details>
