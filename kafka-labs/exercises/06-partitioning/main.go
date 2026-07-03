// Exercise 06 — Partitioning
//
// Goal: understand how a record's KEY selects its partition, and why that
// matters. Kafka only guarantees ordering *within* a partition. The default
// partitioner hashes the record key, so all records with the same key land on
// the same partition — which is how you keep per-key events in order.
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

// produceKeyed produces one record with the given key and value, waits for the
// ack, and returns the partition the broker assigned it. After a synchronous
// produce, the assigned partition is available on the record.
//
// TODO: build a *kgo.Record with Topic/Key/Value, ProduceSync it, and return
// the record's Partition.
func produceKeyed(ctx context.Context, cl *kgo.Client, topic, key, value string) (int32, error) {
	_ = ctx
	_ = cl
	_ = topic
	_ = key
	_ = value
	return -1, errTODO
}

func main() {
	log := labkafka.Logger("06-partitioning")
	ctx, cancel := labkafka.WithTimeout(20e9)
	defer cancel()

	const topic = "labkafka.06-partitioning.demo"
	if err := labkafka.EnsureTopic(ctx, topic, 4, 1); err != nil {
		log.Fatalf("ensure topic: %v", err)
	}

	cl, _ := labkafka.NewClient()
	defer cl.Close()

	// Same key three times -> same partition every time.
	for i := range 3 {
		p, err := produceKeyed(ctx, cl, topic, "user-42", fmt.Sprintf("event-%d", i))
		if err != nil {
			log.Fatalf("produceKeyed: %v", err)
		}
		log.Printf("key=user-42 event-%d -> partition %d", i, p)
	}

	// A different key may map elsewhere.
	p, err := produceKeyed(ctx, cl, topic, "user-99", "hello")
	if err != nil {
		log.Fatalf("produceKeyed: %v", err)
	}
	log.Printf("key=user-99 -> partition %d", p)

	fmt.Println("06-partitioning: done")
}
