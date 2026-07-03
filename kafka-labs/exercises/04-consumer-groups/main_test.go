package main

import (
	"fmt"
	"testing"
	"time"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

// TestGroupCommitResumes verifies the core group guarantee: after one member
// consumes and commits, a fresh member of the same group resumes AFTER the
// committed offset instead of re-reading from the start.
func TestGroupCommitResumes(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}

	// Unique names so reruns don't inherit old committed offsets.
	suffix := time.Now().UnixNano()
	topic := fmt.Sprintf("labkafka.04-groups.test.%d", suffix)
	group := fmt.Sprintf("labkafka.04-groups.test-group.%d", suffix)

	ctx, cancel := labkafka.WithTimeout(30e9)
	defer cancel()
	if err := labkafka.EnsureTopic(ctx, topic, 1, 1); err != nil {
		t.Fatalf("ensure topic: %v", err)
	}
	t.Cleanup(func() {
		c, cancel := labkafka.WithTimeout(15e9)
		defer cancel()
		_ = labkafka.DeleteTopic(c, topic)
	})

	// Seed 4 records on a single partition (stable order).
	want := []string{"m1", "m2", "m3", "m4"}
	producer, _ := labkafka.NewClient()
	for _, v := range want {
		producer.Produce(ctx, &kgo.Record{Topic: topic, Value: []byte(v)}, nil)
	}
	if err := producer.Flush(ctx); err != nil {
		t.Fatalf("flush: %v", err)
	}
	producer.Close()

	newGroupClient := func() *kgo.Client {
		cl, err := labkafka.NewClient(
			kgo.ConsumeTopics(topic),
			kgo.ConsumerGroup(group),
			kgo.DisableAutoCommit(),
			kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		)
		if err != nil {
			t.Fatalf("group client: %v", err)
		}
		return cl
	}

	c1 := newGroupClient()
	first, err := consumeAndCommit(ctx, c1, 2)
	if err != nil {
		c1.Close()
		t.Fatalf("first consumeAndCommit: %v", err)
	}
	c1.Close() // block until the group is left and commits flushed

	if len(first) != 2 || first[0] != "m1" || first[1] != "m2" {
		t.Fatalf("member 1 got %v, want [m1 m2]", first)
	}

	c2 := newGroupClient()
	defer c2.Close()
	second, err := consumeAndCommit(ctx, c2, 2)
	if err != nil {
		t.Fatalf("second consumeAndCommit: %v", err)
	}
	if len(second) != 2 || second[0] != "m3" || second[1] != "m4" {
		t.Fatalf("member 2 got %v, want [m3 m4] (commit did not take effect)", second)
	}
}
