package main

import (
	"testing"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

func TestConsumeN(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}
	const topic = "labkafka.03-consume.test"

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

	// Seed known records.
	want := []string{"one", "two", "three", "four"}
	producer, err := labkafka.NewClient()
	if err != nil {
		t.Fatalf("producer: %v", err)
	}
	for _, v := range want {
		producer.Produce(ctx, &kgo.Record{Topic: topic, Value: []byte(v)}, nil)
	}
	if err := producer.Flush(ctx); err != nil {
		t.Fatalf("flush: %v", err)
	}
	producer.Close()

	cl, err := labkafka.NewClient(
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		t.Fatalf("consumer: %v", err)
	}
	defer cl.Close()

	got, err := consumeN(ctx, cl, len(want))
	if err != nil {
		t.Fatalf("consumeN: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("consumeN returned %d records %v, want %d", len(got), got, len(want))
	}
	// Single partition => order is preserved.
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("record %d = %q, want %q (order should be preserved on one partition)", i, got[i], want[i])
		}
	}
}
