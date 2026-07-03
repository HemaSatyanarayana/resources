// Exercise 02 — Produce
//
// Goal: get records onto a topic and understand the two produce styles:
//   - synchronous: produce and block until the broker acks each record
//   - asynchronous: fire records and handle acks/errors in a callback
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

// produceSync writes each value to topic and blocks until every record has been
// acknowledged by the broker. If any record fails, return the error.
//
// TODO: build a *kgo.Record per value and use cl.ProduceSync.
func produceSync(ctx context.Context, cl *kgo.Client, topic string, values []string) error {
	_ = ctx
	_ = cl
	_ = topic
	_ = values
	return errTODO
}

// produceAsync writes each value to topic using the async cl.Produce API with a
// promise/callback, then waits for all in-flight records to complete before
// returning. Aggregate and return the first error you observe (if any).
//
// TODO: cl.Produce with a callback, and make sure you wait for completion
// (hint: sync.WaitGroup, or cl.Flush).
func produceAsync(ctx context.Context, cl *kgo.Client, topic string, values []string) error {
	_ = ctx
	_ = cl
	_ = topic
	_ = values
	return errTODO
}

func main() {
	log := labkafka.Logger("02-produce")
	ctx, cancel := labkafka.WithTimeout(20e9)
	defer cancel()

	const topic = "labkafka.02-produce.demo"
	if err := labkafka.EnsureTopic(ctx, topic, 1, 1); err != nil {
		log.Fatalf("ensure topic: %v", err)
	}

	cl, err := labkafka.NewClient()
	if err != nil {
		log.Fatalf("client: %v", err)
	}
	defer cl.Close()

	if err := produceSync(ctx, cl, topic, []string{"sync-1", "sync-2", "sync-3"}); err != nil {
		log.Fatalf("produceSync: %v", err)
	}
	log.Printf("produced 3 records synchronously")

	if err := produceAsync(ctx, cl, topic, []string{"async-1", "async-2", "async-3"}); err != nil {
		log.Fatalf("produceAsync: %v", err)
	}
	log.Printf("produced 3 records asynchronously")

	fmt.Println("02-produce: done")
}
