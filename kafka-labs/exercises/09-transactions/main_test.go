package main

import (
	"fmt"
	"testing"
	"time"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

// TestAbortHidesRecords is the core transaction spec: aborted records stay
// invisible to a ReadCommitted consumer, while committed records show up.
func TestAbortHidesRecords(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}

	suffix := time.Now().UnixNano()
	topic := fmt.Sprintf("labkafka.09-transactions.test.%d", suffix)
	txnID := fmt.Sprintf("labkafka.09-txn.test.%d", suffix)

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

	cl, err := newTxnProducer(txnID)
	if err != nil {
		t.Fatalf("newTxnProducer: %v", err)
	}
	defer cl.Close()

	if err := produceTxn(ctx, cl, topic, []string{"ghost-1", "ghost-2"}, false); err != nil {
		t.Fatalf("aborted txn: %v", err)
	}
	if err := produceTxn(ctx, cl, topic, []string{"real-1", "real-2"}, true); err != nil {
		t.Fatalf("committed txn: %v", err)
	}

	// ReadCommitted consumer must see ONLY the committed records.
	consumer, err := labkafka.NewClient(
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.FetchIsolationLevel(kgo.ReadCommitted()),
	)
	if err != nil {
		t.Fatalf("consumer: %v", err)
	}
	defer consumer.Close()

	got := map[string]bool{}
	// Poll a few rounds; committed set is small.
	deadline := time.Now().Add(10 * time.Second)
	for len(got) < 2 && time.Now().Before(deadline) {
		fs := consumer.PollFetches(ctx)
		if errs := fs.Errors(); len(errs) > 0 {
			t.Fatalf("poll: %v", errs)
		}
		fs.EachRecord(func(r *kgo.Record) { got[string(r.Value)] = true })
	}

	if got["ghost-1"] || got["ghost-2"] {
		t.Fatalf("aborted records became visible: %v", got)
	}
	if !got["real-1"] || !got["real-2"] {
		t.Fatalf("committed records missing: %v", got)
	}
}
