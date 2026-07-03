// Exercise 05 — Offsets
//
// Goal: work with offsets directly. Offsets are how Kafka tracks position. In
// this exercise you:
//   - read a partition's end offset (the "log end offset") with kadm
//   - seek a consumer to an arbitrary offset and read from there
//
// This is the machinery underneath consumer groups — useful for replay,
// rewinding, or skipping ahead.
//
// Fill in the TODOs, then `go test ./...` and `go run .`.
package main

import (
	"context"
	"errors"
	"fmt"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

var errTODO = errors.New("TODO: not implemented")

// endOffset returns the log end offset (one past the last record) for
// topic/partition. For an empty topic this is 0; after producing N records to a
// single partition it is N.
//
// TODO: use adm.ListEndOffsets and pull out the offset for topic/partition.
func endOffset(ctx context.Context, adm *kadm.Client, topic string, partition int32) (int64, error) {
	_ = ctx
	_ = adm
	_ = topic
	_ = partition
	return 0, errTODO
}

// readFrom consumes up to n record values from topic/partition starting exactly
// at startOffset. The client is created fresh here (no group) so you control the
// start position precisely.
//
// TODO: build a client with kgo.ConsumePartitions(topic, {partition: offset}),
// where offset is kgo.NewOffset().At(startOffset); then poll n records.
func readFrom(ctx context.Context, topic string, partition int32, startOffset int64, n int) ([]string, error) {
	_ = ctx
	_ = topic
	_ = partition
	_ = startOffset
	_ = n
	return nil, errTODO
}

func main() {
	log := labkafka.Logger("05-offsets")
	ctx, cancel := labkafka.WithTimeout(20e9)
	defer cancel()

	const topic = "labkafka.05-offsets.demo"
	if err := labkafka.EnsureTopic(ctx, topic, 1, 1); err != nil {
		log.Fatalf("ensure topic: %v", err)
	}

	// Seed 5 records.
	producer, _ := labkafka.NewClient()
	for i := range 5 {
		producer.Produce(ctx, &kgo.Record{Topic: topic, Value: fmt.Appendf(nil, "r%d", i)}, nil)
	}
	_ = producer.Flush(ctx)
	producer.Close()

	adm, _ := labkafka.NewAdmin()
	defer adm.Close()

	end, err := endOffset(ctx, adm, topic, 0)
	if err != nil {
		log.Fatalf("endOffset: %v", err)
	}
	log.Printf("end offset for %q p0 = %d", topic, end)

	// Read the last two records by seeking to end-2.
	tail, err := readFrom(ctx, topic, 0, end-2, 2)
	if err != nil {
		log.Fatalf("readFrom: %v", err)
	}
	log.Printf("last 2 records (from offset %d): %v", end-2, tail)

	fmt.Println("05-offsets: done")
}
