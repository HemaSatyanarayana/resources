package main

import (
	"testing"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

// TestEncodeDecodeUnit checks the round-trip in isolation — no broker needed,
// so it always runs.
func TestEncodeDecodeUnit(t *testing.T) {
	in := Order{ID: "o-1", Item: "bolt", Qty: 12, Amount: 4.50}
	b, err := encodeOrder(in)
	if err != nil {
		t.Fatalf("encodeOrder: %v", err)
	}
	out, err := decodeOrder(b)
	if err != nil {
		t.Fatalf("decodeOrder: %v", err)
	}
	if out != in {
		t.Fatalf("round-trip mismatch: got %+v, want %+v", out, in)
	}
}

// TestOrderOverKafka sends an encoded order through the broker and decodes it on
// the other side.
func TestOrderOverKafka(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}
	const topic = "labkafka.07-serialization.test"

	ctx, cancel := labkafka.WithTimeout(20e9)
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

	want := Order{ID: "o-777", Item: "gizmo", Qty: 2, Amount: 19.99}
	value, err := encodeOrder(want)
	if err != nil {
		t.Fatalf("encodeOrder: %v", err)
	}

	producer, _ := labkafka.NewClient()
	if err := producer.ProduceSync(ctx, &kgo.Record{Topic: topic, Value: value}).FirstErr(); err != nil {
		t.Fatalf("produce: %v", err)
	}
	producer.Close()

	cl, _ := labkafka.NewClient(
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	defer cl.Close()

	fs := cl.PollFetches(ctx)
	if errs := fs.Errors(); len(errs) > 0 {
		t.Fatalf("poll: %v", errs)
	}
	recs := fs.Records()
	if len(recs) == 0 {
		t.Fatalf("no records consumed")
	}
	got, err := decodeOrder(recs[0].Value)
	if err != nil {
		t.Fatalf("decodeOrder: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
