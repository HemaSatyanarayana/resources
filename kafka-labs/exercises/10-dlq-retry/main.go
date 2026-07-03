// Exercise 10 — Dead-letter queue & retry
//
// Goal: handle "poison" records that fail processing without blocking the whole
// consumer. The standard pattern: try to process each record; if it fails,
// publish it to a dead-letter topic (DLQ) with the failure reason attached, then
// move on. A separate process can inspect or replay the DLQ later.
//
// isPoison below is the (provided) business rule that decides success/failure.
// Your job is the routing.
//
// Fill in the TODOs, then `go test ./...` and `go run .`.
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

var errTODO = errors.New("TODO: not implemented")

// isPoison is the provided processing rule: any value containing "poison"
// fails. (In a real app this would be your handler returning an error.)
func isPoison(value string) bool {
	return strings.Contains(value, "poison")
}

// routeToDLQ publishes a failed record to dlqTopic, preserving its key and
// value and attaching the failure reason as a record header named "error".
// Wait for the ack before returning.
//
// TODO: build a *kgo.Record for dlqTopic with r.Key / r.Value and
// Headers: []kgo.RecordHeader{{Key: "error", Value: []byte(cause.Error())}};
// ProduceSync it and return the first error.
func routeToDLQ(ctx context.Context, producer *kgo.Client, dlqTopic string, r *kgo.Record, cause error) error {
	_ = ctx
	_ = producer
	_ = dlqTopic
	_ = r
	_ = cause
	return errTODO
}

// processBatch consumes up to n records. For each, it runs the processing rule:
// good records go into `ok`; failing records are routed to the DLQ via
// routeToDLQ and their values collected into `dead`. Return once n records have
// been handled (or ctx is done).
//
// TODO: loop consumer.PollFetches(ctx); for each record, if isPoison(value)
// call routeToDLQ and append to dead, else append to ok; stop after n handled.
func processBatch(ctx context.Context, consumer, dlqProducer *kgo.Client, dlqTopic string, n int) (ok, dead []string, err error) {
	_ = ctx
	_ = consumer
	_ = dlqProducer
	_ = dlqTopic
	_ = n
	return nil, nil, errTODO
}

func main() {
	log := labkafka.Logger("10-dlq-retry")
	ctx, cancel := labkafka.WithTimeout(30e9)
	defer cancel()

	const topic = "labkafka.10-dlq.demo"
	const dlqTopic = "labkafka.10-dlq.demo.DLQ"
	for _, tp := range []string{topic, dlqTopic} {
		if err := labkafka.EnsureTopic(ctx, tp, 1, 1); err != nil {
			log.Fatalf("ensure topic %q: %v", tp, err)
		}
	}

	// Seed a mix of good and poison records.
	producer, _ := labkafka.NewClient()
	for _, v := range []string{"ok-1", "poison-A", "ok-2", "poison-B", "ok-3"} {
		producer.Produce(ctx, &kgo.Record{Topic: topic, Value: []byte(v)}, nil)
	}
	_ = producer.Flush(ctx)

	consumer, _ := labkafka.NewClient(
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	defer consumer.Close()

	ok, dead, err := processBatch(ctx, consumer, producer, dlqTopic, 5)
	if err != nil {
		log.Fatalf("processBatch: %v", err)
	}
	producer.Close()
	log.Printf("processed OK: %v", ok)
	log.Printf("routed to DLQ: %v", dead)

	fmt.Println("10-dlq-retry: done")
}
