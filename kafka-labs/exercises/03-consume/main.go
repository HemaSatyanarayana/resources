// Exercise 03 — Consume
//
// Goal: read records from a topic without a consumer group (a "simple"
// consumer). You configure which topic(s) and start offset on the client, then
// poll fetches in a loop and iterate the records.
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

// consumeN polls the client until it has collected n record values (or ctx is
// done) and returns them in the order received. The client is expected to have
// been created with kgo.ConsumeTopics(...) already.
//
// TODO: loop on cl.PollFetches(ctx); check fs.Errors(); collect r.Value via
// fs.EachRecord until you have n (or ctx.Err() != nil).
func consumeN(ctx context.Context, cl *kgo.Client, n int) ([]string, error) {
	_ = ctx
	_ = cl
	_ = n
	return nil, errTODO
}

func main() {
	log := labkafka.Logger("03-consume")
	ctx, cancel := labkafka.WithTimeout(20e9)
	defer cancel()

	const topic = "labkafka.03-consume.demo"
	if err := labkafka.EnsureTopic(ctx, topic, 1, 1); err != nil {
		log.Fatalf("ensure topic: %v", err)
	}

	// Seed a few records so there is something to read.
	producer, err := labkafka.NewClient()
	if err != nil {
		log.Fatalf("producer: %v", err)
	}
	for _, v := range []string{"hello", "kafka", "world"} {
		producer.Produce(ctx, &kgo.Record{Topic: topic, Value: []byte(v)}, nil)
	}
	if err := producer.Flush(ctx); err != nil {
		log.Fatalf("flush: %v", err)
	}
	producer.Close()

	// Consume from the start of the topic.
	cl, err := labkafka.NewClient(
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		log.Fatalf("consumer: %v", err)
	}
	defer cl.Close()

	values, err := consumeN(ctx, cl, 3)
	if err != nil {
		log.Fatalf("consumeN: %v", err)
	}
	log.Printf("consumed %d records: %v", len(values), values)

	fmt.Println("03-consume: done")
}
