// Exercise 07 — Serialization
//
// Goal: Kafka records are just bytes. Real systems put structured data on the
// wire, so you need to encode on produce and decode on consume. Here we use
// JSON; the same shape applies to Avro/Protobuf (with a schema registry).
//
// Fill in the TODOs, then `go test ./...` and `go run .`.
package main

import (
	"errors"
	"fmt"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

var errTODO = errors.New("TODO: not implemented")

// Order is the domain event we put on the topic.
type Order struct {
	ID     string  `json:"id"`
	Item   string  `json:"item"`
	Qty    int     `json:"qty"`
	Amount float64 `json:"amount"`
}

// encodeOrder serializes o to the bytes stored in a record's Value.
//
// TODO: marshal o to JSON.
func encodeOrder(o Order) ([]byte, error) {
	_ = o
	return nil, errTODO
}

// decodeOrder is the inverse of encodeOrder.
//
// TODO: unmarshal b into an Order.
func decodeOrder(b []byte) (Order, error) {
	_ = b
	return Order{}, errTODO
}

func main() {
	log := labkafka.Logger("07-serialization")
	ctx, cancel := labkafka.WithTimeout(20e9)
	defer cancel()

	const topic = "labkafka.07-serialization.demo"
	if err := labkafka.EnsureTopic(ctx, topic, 1, 1); err != nil {
		log.Fatalf("ensure topic: %v", err)
	}

	order := Order{ID: "o-1001", Item: "widget", Qty: 3, Amount: 29.97}

	// Produce the encoded order.
	value, err := encodeOrder(order)
	if err != nil {
		log.Fatalf("encodeOrder: %v", err)
	}
	producer, _ := labkafka.NewClient()
	if err := producer.ProduceSync(ctx, &kgo.Record{Topic: topic, Key: []byte(order.ID), Value: value}).FirstErr(); err != nil {
		log.Fatalf("produce: %v", err)
	}
	producer.Close()
	log.Printf("produced order %+v", order)

	// Consume and decode it back.
	cl, _ := labkafka.NewClient(
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	defer cl.Close()
	fs := cl.PollFetches(ctx)
	if errs := fs.Errors(); len(errs) > 0 {
		log.Fatalf("poll: %v", errs)
	}
	rec := fs.Records()[0]
	got, err := decodeOrder(rec.Value)
	if err != nil {
		log.Fatalf("decodeOrder: %v", err)
	}
	log.Printf("decoded order %+v", got)

	fmt.Println("07-serialization: done")
}
