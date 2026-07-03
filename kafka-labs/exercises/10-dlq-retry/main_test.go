package main

import (
	"testing"

	"kafka-labs/internal/labkafka"

	"github.com/twmb/franz-go/pkg/kgo"
)

// TestDLQRouting seeds good and poison records, then checks that good ones are
// processed and poison ones land in the DLQ with an "error" header.
func TestDLQRouting(t *testing.T) {
	if !labkafka.BrokersConfigured() {
		t.Skip("set KAFKA_BROKERS to run this integration test")
	}
	const topic = "labkafka.10-dlq.test"
	const dlqTopic = "labkafka.10-dlq.test.DLQ"

	ctx, cancel := labkafka.WithTimeout(30e9)
	defer cancel()
	for _, tp := range []string{topic, dlqTopic} {
		_ = labkafka.DeleteTopic(ctx, tp)
		if err := labkafka.EnsureTopic(ctx, tp, 1, 1); err != nil {
			t.Fatalf("ensure topic %q: %v", tp, err)
		}
	}
	t.Cleanup(func() {
		c, cancel := labkafka.WithTimeout(15e9)
		defer cancel()
		_ = labkafka.DeleteTopic(c, topic)
		_ = labkafka.DeleteTopic(c, dlqTopic)
	})

	producer, _ := labkafka.NewClient()
	seed := []string{"ok-1", "poison-A", "ok-2", "poison-B", "ok-3"}
	for _, v := range seed {
		producer.Produce(ctx, &kgo.Record{Topic: topic, Value: []byte(v)}, nil)
	}
	if err := producer.Flush(ctx); err != nil {
		t.Fatalf("flush: %v", err)
	}

	consumer, _ := labkafka.NewClient(
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	defer consumer.Close()

	ok, dead, err := processBatch(ctx, consumer, producer, dlqTopic, len(seed))
	if err != nil {
		t.Fatalf("processBatch: %v", err)
	}
	producer.Close()

	assertSet(t, "processed", ok, []string{"ok-1", "ok-2", "ok-3"})
	assertSet(t, "dead", dead, []string{"poison-A", "poison-B"})

	// Verify the DLQ topic really received the poison records, with the header.
	dlqConsumer, _ := labkafka.NewClient(
		kgo.ConsumeTopics(dlqTopic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	defer dlqConsumer.Close()

	got := map[string]string{} // value -> error header
	for len(got) < 2 && ctx.Err() == nil {
		fs := dlqConsumer.PollFetches(ctx)
		if errs := fs.Errors(); len(errs) > 0 {
			t.Fatalf("dlq poll: %v", errs)
		}
		fs.EachRecord(func(r *kgo.Record) {
			var hdr string
			for _, h := range r.Headers {
				if h.Key == "error" {
					hdr = string(h.Value)
				}
			}
			got[string(r.Value)] = hdr
		})
	}

	for _, v := range []string{"poison-A", "poison-B"} {
		hdr, ok := got[v]
		if !ok {
			t.Fatalf("DLQ missing record %q", v)
		}
		if hdr == "" {
			t.Fatalf("DLQ record %q has no non-empty \"error\" header", v)
		}
	}
}

func assertSet(t *testing.T, label string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s = %v (len %d), want %v (len %d)", label, got, len(got), want, len(want))
	}
	seen := map[string]bool{}
	for _, g := range got {
		seen[g] = true
	}
	for _, w := range want {
		if !seen[w] {
			t.Fatalf("%s missing %q; got %v", label, w, got)
		}
	}
}
