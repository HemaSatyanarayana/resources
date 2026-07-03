package main

import (
	"fmt"
	"testing"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

func TestEndOffsetAndSeek(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}
	const topic = "labkafka.05-offsets.test"

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

	const total = 5
	producer, _ := labkafka.NewClient()
	for i := range total {
		producer.Produce(ctx, &kgo.Record{Topic: topic, Value: fmt.Appendf(nil, "r%d", i)}, nil)
	}
	if err := producer.Flush(ctx); err != nil {
		t.Fatalf("flush: %v", err)
	}
	producer.Close()

	adm, _ := labkafka.NewAdmin()
	defer adm.Close()

	end, err := endOffset(ctx, adm, topic, 0)
	if err != nil {
		t.Fatalf("endOffset: %v", err)
	}
	if end != total {
		t.Fatalf("endOffset = %d, want %d", end, total)
	}

	// Seek to offset 2 and read the tail: r2, r3, r4.
	got, err := readFrom(ctx, topic, 0, 2, 3)
	if err != nil {
		t.Fatalf("readFrom: %v", err)
	}
	want := []string{"r2", "r3", "r4"}
	if len(got) != len(want) {
		t.Fatalf("readFrom returned %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("record %d = %q, want %q", i, got[i], want[i])
		}
	}
}
