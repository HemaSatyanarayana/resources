package main

import (
	"testing"

	"kafka-labs/internal/labkafka"
)

// TestSameKeySamePartition verifies the key ordering guarantee: identical keys
// always resolve to one partition, so their records stay ordered.
func TestSameKeySamePartition(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}
	const topic = "labkafka.06-partitioning.test"

	ctx, cancel := labkafka.WithTimeout(20e9)
	defer cancel()
	_ = labkafka.DeleteTopic(ctx, topic)
	if err := labkafka.EnsureTopic(ctx, topic, 4, 1); err != nil {
		t.Fatalf("ensure topic: %v", err)
	}
	t.Cleanup(func() {
		c, cancel := labkafka.WithTimeout(15e9)
		defer cancel()
		_ = labkafka.DeleteTopic(c, topic)
	})

	cl, err := labkafka.NewClient()
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	defer cl.Close()

	// The same key must always land on the same partition.
	var first int32 = -1
	for i := range 5 {
		p, err := produceKeyed(ctx, cl, topic, "steady-key", "v")
		if err != nil {
			t.Fatalf("produceKeyed[%d]: %v", i, err)
		}
		if p < 0 {
			t.Fatalf("produceKeyed returned partition %d (did you return r.Partition?)", p)
		}
		if first == -1 {
			first = p
		} else if p != first {
			t.Fatalf("same key mapped to partition %d then %d; keys must be stable", first, p)
		}
	}
}
