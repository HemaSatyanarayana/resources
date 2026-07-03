// Exercise 04 — Consumer groups
//
// Goal: consume as a member of a consumer group and commit offsets so progress
// survives a restart. A group lets multiple consumers share a topic's
// partitions, and the broker remembers each group's committed offset.
//
// Here we do MANUAL commits: poll a batch, then commit exactly what we polled.
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

// consumeAndCommit polls until it has n records (or ctx is done), then commits
// the offsets for exactly those records so a later consumer in the same group
// resumes after them. Return the record values in order.
//
// The client is created with kgo.ConsumerGroup(...) and kgo.DisableAutoCommit()
// by the caller, so committing is your responsibility.
//
// TODO: loop PollFetches; collect values; keep the kgo.Fetches (or records)
// so you can call cl.CommitRecords(ctx, recs...) once you have enough.
func consumeAndCommit(ctx context.Context, cl *kgo.Client, n int) ([]string, error) {
	_ = ctx
	_ = cl
	_ = n
	return nil, errTODO
}

func main() {
	log := labkafka.Logger("04-consumer-groups")
	ctx, cancel := labkafka.WithTimeout(30e9)
	defer cancel()

	const topic = "labkafka.04-groups.demo"
	const group = "labkafka.04-groups.demo-group"
	if err := labkafka.EnsureTopic(ctx, topic, 1, 1); err != nil {
		log.Fatalf("ensure topic: %v", err)
	}

	// Seed records.
	producer, _ := labkafka.NewClient()
	for _, v := range []string{"m1", "m2", "m3", "m4"} {
		producer.Produce(ctx, &kgo.Record{Topic: topic, Value: []byte(v)}, nil)
	}
	_ = producer.Flush(ctx)
	producer.Close()

	newGroupClient := func() *kgo.Client {
		cl, err := labkafka.NewClient(
			kgo.ConsumeTopics(topic),
			kgo.ConsumerGroup(group),
			kgo.DisableAutoCommit(),
			kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		)
		if err != nil {
			log.Fatalf("group client: %v", err)
		}
		return cl
	}

	// First member reads two records and commits.
	c1 := newGroupClient()
	first, err := consumeAndCommit(ctx, c1, 2)
	if err != nil {
		log.Fatalf("first consume: %v", err)
	}
	c1.Close() // triggers a final commit + leaves the group
	log.Printf("member 1 consumed+committed: %v", first)

	// A fresh member of the SAME group should resume AFTER the committed offset.
	c2 := newGroupClient()
	defer c2.Close()
	second, err := consumeAndCommit(ctx, c2, 2)
	if err != nil {
		log.Fatalf("second consume: %v", err)
	}
	log.Printf("member 2 resumed at committed offset: %v", second)

	fmt.Println("04-consumer-groups: done")
}
