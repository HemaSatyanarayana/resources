// Exercise 08 — Idempotent delivery
//
// Goal: configure a producer so that internal retries never create duplicate
// records. Kafka's idempotent producer tags each record with a producer id and
// sequence number, so the broker can drop a retried record it has already seen.
//
// franz-go enables idempotency by DEFAULT and, to keep it, requires acks=all
// (AllISRAcks). The mistake to avoid is weakening acks or disabling idempotency
// for "speed" and silently getting duplicates on retry.
//
// Fill in the TODOs, then `go test ./...` and `go run .`.
package main

import (
	"context"
	"errors"
	"fmt"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

var errTODO = errors.New("TODO: not implemented")

// newIdempotentProducer returns a client configured for safe, idempotent
// delivery: acks=all and idempotency left enabled. Do NOT pass
// kgo.DisableIdempotentWrite.
//
// TODO: labkafka.NewClient(kgo.RequiredAcks(kgo.AllISRAcks())) and return it.
func newIdempotentProducer() (*kgo.Client, error) {
	return nil, errTODO
}

// produceBatch produces every value to topic and blocks until all are acked.
// With the idempotent producer above, retried records are de-duplicated by the
// broker, so exactly len(values) records are appended.
//
// TODO: ProduceSync all records (or Produce + Flush) and return the first error.
func produceBatch(ctx context.Context, cl *kgo.Client, topic string, values []string) error {
	_ = ctx
	_ = cl
	_ = topic
	_ = values
	return errTODO
}

func main() {
	log := labkafka.Logger("08-idempotent-delivery")
	ctx, cancel := labkafka.WithTimeout(20e9)
	defer cancel()

	const topic = "labkafka.08-idempotent.demo"
	if err := labkafka.EnsureTopic(ctx, topic, 1, 1); err != nil {
		log.Fatalf("ensure topic: %v", err)
	}

	cl, err := newIdempotentProducer()
	if err != nil {
		log.Fatalf("newIdempotentProducer: %v", err)
	}
	defer cl.Close()

	values := []string{"e1", "e2", "e3", "e4", "e5"}
	if err := produceBatch(ctx, cl, topic, values); err != nil {
		log.Fatalf("produceBatch: %v", err)
	}
	log.Printf("produced %d records with idempotent delivery", len(values))

	fmt.Println("08-idempotent-delivery: done")
}
