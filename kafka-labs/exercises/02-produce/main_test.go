package main

import (
	"testing"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

// consumeValues reads up to n record values from topic starting at the earliest
// offset. It is test-only scaffolding so we can verify what was produced.
func consumeValues(t *testing.T, topic string, n int) []string {
	t.Helper()
	cl, err := labkafka.NewClient(
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		t.Fatalf("consumer client: %v", err)
	}
	defer cl.Close()

	ctx, cancel := labkafka.WithTimeout(15e9)
	defer cancel()

	var out []string
	for len(out) < n {
		fs := cl.PollFetches(ctx)
		if errs := fs.Errors(); len(errs) > 0 {
			t.Fatalf("poll: %v", errs)
		}
		fs.EachRecord(func(r *kgo.Record) { out = append(out, string(r.Value)) })
		if ctx.Err() != nil {
			break
		}
	}
	return out
}

func TestProduceSync(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}
	const topic = "labkafka.02-produce.sync.test"
	seedTopic(t, topic)

	cl := newClient(t)
	defer cl.Close()

	want := []string{"a", "b", "c"}
	if err := produceSync(t.Context(), cl, topic, want); err != nil {
		t.Fatalf("produceSync: %v", err)
	}

	got := consumeValues(t, topic, len(want))
	assertSameSet(t, got, want)
}

func TestProduceAsync(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}
	const topic = "labkafka.02-produce.async.test"
	seedTopic(t, topic)

	cl := newClient(t)
	defer cl.Close()

	want := []string{"x", "y", "z"}
	if err := produceAsync(t.Context(), cl, topic, want); err != nil {
		t.Fatalf("produceAsync: %v", err)
	}

	got := consumeValues(t, topic, len(want))
	assertSameSet(t, got, want)
}

// --- small shared test helpers ---

func newClient(t *testing.T) *kgo.Client {
	t.Helper()
	cl, err := labkafka.NewClient()
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	return cl
}

func seedTopic(t *testing.T, topic string) {
	t.Helper()
	ctx, cancel := labkafka.WithTimeout(15e9)
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
}

func assertSameSet(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d records %v, want %d %v", len(got), got, len(want), want)
	}
	seen := make(map[string]int, len(got))
	for _, g := range got {
		seen[g]++
	}
	for _, w := range want {
		if seen[w] == 0 {
			t.Fatalf("missing expected value %q in %v", w, got)
		}
		seen[w]--
	}
}
