// Exercise 01 — Admin
//
// Goal: manage topics directly with the kadm admin client — create, list,
// describe, and delete. This is the foundation every other exercise builds on.
//
// Fill in the TODOs below. Run `go test ./...` in this directory to check your
// work, then `go run .` to watch it end to end.
package main

import (
	"context"
	"errors"
	"fmt"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kadm"
)

// errTODO marks an unfinished exercise step. Replace the returns below with real
// implementations and this sentinel disappears.
var errTODO = errors.New("TODO: not implemented")

// createTopic creates a topic with the given partition count and replication
// factor. It must be idempotent: if the topic already exists, return nil rather
// than an error (hint: inspect the per-topic response and compare the Kafka
// error code to kerr.TopicAlreadyExists).
//
// TODO: call adm.CreateTopic and handle the "already exists" case.
func createTopic(ctx context.Context, adm *kadm.Client, topic string, partitions int32, rf int16) error {
	_ = ctx
	_ = adm
	_ = topic
	_ = partitions
	_ = rf
	return errTODO
}

// topicExists reports whether topic is present on the cluster.
//
// TODO: use adm.ListTopics (or adm.TopicDetails) and check membership.
func topicExists(ctx context.Context, adm *kadm.Client, topic string) (bool, error) {
	_ = ctx
	_ = adm
	_ = topic
	return false, errTODO
}

// partitionCount returns how many partitions topic has.
//
// TODO: describe the topic and count its partitions.
func partitionCount(ctx context.Context, adm *kadm.Client, topic string) (int, error) {
	_ = ctx
	_ = adm
	_ = topic
	return 0, errTODO
}

// deleteTopic removes topic. Treat "unknown topic" as success so it is safe to
// call during cleanup.
//
// TODO: call adm.DeleteTopics and handle kerr.UnknownTopicOrPartition.
func deleteTopic(ctx context.Context, adm *kadm.Client, topic string) error {
	_ = ctx
	_ = adm
	_ = topic
	return errTODO
}

func main() {
	log := labkafka.Logger("01-admin")
	ctx, cancel := labkafka.WithTimeout(15e9) // 15s
	defer cancel()

	adm, err := labkafka.NewAdmin()
	if err != nil {
		log.Fatalf("admin client: %v", err)
	}
	defer adm.Close()

	const topic = "labkafka.01-admin.demo"

	if err := createTopic(ctx, adm, topic, 3, 1); err != nil {
		log.Fatalf("createTopic: %v", err)
	}
	log.Printf("created (or already had) topic %q", topic)

	exists, err := topicExists(ctx, adm, topic)
	if err != nil {
		log.Fatalf("topicExists: %v", err)
	}
	log.Printf("topicExists(%q) = %v", topic, exists)

	n, err := partitionCount(ctx, adm, topic)
	if err != nil {
		log.Fatalf("partitionCount: %v", err)
	}
	log.Printf("%q has %d partitions", topic, n)

	if err := deleteTopic(ctx, adm, topic); err != nil {
		log.Fatalf("deleteTopic: %v", err)
	}
	log.Printf("deleted topic %q", topic)

	fmt.Println("01-admin: done")
}
