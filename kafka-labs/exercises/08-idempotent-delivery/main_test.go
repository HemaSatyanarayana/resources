package main

import (
	"testing"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

// TestIdempotentNoDuplicates produces a known set and verifies the topic
// contains exactly those records — no duplicates, none missing.
func TestIdempotentNoDuplicates(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}
	const topic = "labkafka.08-idempotent.test"

	ctx, cancel := labkafka.WithTimeout(25e9)
	defer cancel()
	_ = labkafka.DeleteTopic(ctx, topic)
	if err := labkafka.EnsureTopic(ctx, topic, 1, 1); err != nil {
		t.Fatalf("ensure topic: %v", err)
	}
	t.Cleanup(func() {
		c, cancel := labkafka.WithTimeout(15e9)
		defer cancel()
		_ = labkafka.DeleteTopic(c, topic)
	})

	cl, err := newIdempotentProducer()
	if err != nil {
		t.Fatalf("newIdempotentProducer: %v", err)
	}
	defer cl.Close()

	want := []string{"e1", "e2", "e3", "e4", "e5"}
	if err := produceBatch(ctx, cl, topic, want); err != nil {
		t.Fatalf("produceBatch: %v", err)
	}

	// Read everything back and count each value.
	consumer, _ := labkafka.NewClient(
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	defer consumer.Close()

	counts := map[string]int{}
	total := 0
	for total < len(want) {
		fs := consumer.PollFetches(ctx)
		if errs := fs.Errors(); len(errs) > 0 {
			t.Fatalf("poll: %v", errs)
		}
		fs.EachRecord(func(r *kgo.Record) {
			counts[string(r.Value)]++
			total++
		})
		if ctx.Err() != nil {
			break
		}
	}

	if total != len(want) {
		t.Fatalf("consumed %d records, want exactly %d (%v)", total, len(want), counts)
	}
	for _, v := range want {
		if counts[v] != 1 {
			t.Fatalf("value %q appeared %d times, want exactly 1", v, counts[v])
		}
	}
}
