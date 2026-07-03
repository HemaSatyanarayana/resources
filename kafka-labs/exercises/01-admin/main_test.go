package main

import (
	"testing"

	"kafka-labs/internal/labkafka"
)

// TestAdminLifecycle is the spec for exercise 01. It exercises the full
// create → exists → describe → delete lifecycle against a live broker.
// It skips when KAFKA_BROKERS is unset so the lab is green without Kafka.
func TestAdminLifecycle(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}

	const topic = "labkafka.01-admin.test"
	const wantPartitions = 3

	adm, err := labkafka.NewAdmin()
	if err != nil {
		t.Fatalf("new admin: %v", err)
	}
	defer adm.Close()

	ctx, cancel := labkafka.WithTimeout(20e9)
	defer cancel()

	// Clean slate, and clean up afterwards.
	_ = deleteTopic(ctx, adm, topic)
	t.Cleanup(func() { _ = deleteTopic(ctx, adm, topic) })

	if err := createTopic(ctx, adm, topic, wantPartitions, 1); err != nil {
		t.Fatalf("createTopic: %v", err)
	}

	// Idempotency: a second create must not error.
	if err := createTopic(ctx, adm, topic, wantPartitions, 1); err != nil {
		t.Fatalf("createTopic (second call must be idempotent): %v", err)
	}

	exists, err := topicExists(ctx, adm, topic)
	if err != nil {
		t.Fatalf("topicExists: %v", err)
	}
	if !exists {
		t.Fatalf("topicExists = false, want true after create")
	}

	n, err := partitionCount(ctx, adm, topic)
	if err != nil {
		t.Fatalf("partitionCount: %v", err)
	}
	if n != wantPartitions {
		t.Fatalf("partitionCount = %d, want %d", n, wantPartitions)
	}

	if err := deleteTopic(ctx, adm, topic); err != nil {
		t.Fatalf("deleteTopic: %v", err)
	}

	exists, err = topicExists(ctx, adm, topic)
	if err != nil {
		t.Fatalf("topicExists after delete: %v", err)
	}
	if exists {
		t.Fatalf("topicExists = true after delete, want false")
	}
}
