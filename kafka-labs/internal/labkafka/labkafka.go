// Package labkafka holds the shared plumbing for every exercise in this lab.
//
// The exercises focus on ONE Kafka concept each, so the boring-but-necessary
// bits (reading the broker list, building a client, creating/cleaning up
// topics, logging) live here and are fully implemented. You should not need to
// edit this package to complete an exercise — you only call into it.
//
// Broker addresses come from the KAFKA_BROKERS environment variable
// (comma-separated). If it is unset we fall back to localhost:9092, which is
// what the lab's local broker listens on.
package labkafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

// isKafkaErr reports whether err carries the given Kafka error code. kadm
// surfaces broker errors as *kerr.Error; we compare by numeric code rather than
// pointer identity so wrapping never breaks the check.
func isKafkaErr(err error, code int16) bool {
	var ke *kerr.Error
	return errors.As(err, &ke) && ke.Code == code
}

// DefaultBroker is used when KAFKA_BROKERS is not set.
const DefaultBroker = "localhost:9092"

// brokersEnv is the environment variable the whole lab reads for its broker
// list. Tests also key off it: when it is unset, integration tests skip.
const brokersEnv = "KAFKA_BROKERS"

// Brokers returns the configured broker list as a slice, splitting
// KAFKA_BROKERS on commas and trimming whitespace. It always returns at least
// one entry (DefaultBroker) so callers never have to nil-check.
func Brokers() []string {
	raw := os.Getenv(brokersEnv)
	if strings.TrimSpace(raw) == "" {
		return []string{DefaultBroker}
	}
	var out []string
	for b := range strings.SplitSeq(raw, ",") {
		if b = strings.TrimSpace(b); b != "" {
			out = append(out, b)
		}
	}
	if len(out) == 0 {
		return []string{DefaultBroker}
	}
	return out
}

// BrokersConfigured reports whether KAFKA_BROKERS was explicitly set. Tests use
// this to decide whether a live broker is expected; when false they skip rather
// than fail, so `go test ./...` is green on a machine with no Kafka.
func BrokersConfigured() bool {
	return strings.TrimSpace(os.Getenv(brokersEnv)) != ""
}

// NewClient builds a franz-go client pointed at the lab's brokers. Any extra
// options you pass are appended after the defaults, so you can override them
// (e.g. kgo.ConsumerGroup, kgo.ConsumeTopics) from inside an exercise.
func NewClient(opts ...kgo.Opt) (*kgo.Client, error) {
	base := []kgo.Opt{
		kgo.SeedBrokers(Brokers()...),
		kgo.ClientID("kafka-labs"),
	}
	return kgo.NewClient(append(base, opts...)...)
}

// NewAdmin builds a kadm admin client for topic/partition/offset management.
// The caller owns the returned client and must Close it.
func NewAdmin() (*kadm.Client, error) {
	cl, err := NewClient()
	if err != nil {
		return nil, err
	}
	return kadm.NewClient(cl), nil
}

// EnsureTopic creates topic if it does not already exist. It is idempotent: a
// "topic already exists" response is treated as success, so exercises can call
// it on every run without caring about prior state.
func EnsureTopic(ctx context.Context, topic string, partitions int32, replicationFactor int16) error {
	adm, err := NewAdmin()
	if err != nil {
		return err
	}
	defer adm.Close()

	resp, err := adm.CreateTopic(ctx, partitions, replicationFactor, nil, topic)
	if err != nil {
		return fmt.Errorf("create topic %q: %w", topic, err)
	}
	if resp.Err != nil && !isKafkaErr(resp.Err, kerr.TopicAlreadyExists.Code) {
		return fmt.Errorf("create topic %q: %w", topic, resp.Err)
	}
	return nil
}

// DeleteTopic removes topic, ignoring "unknown topic" so cleanup is safe to run
// even if the topic was never created. Handy in test teardown.
func DeleteTopic(ctx context.Context, topic string) error {
	adm, err := NewAdmin()
	if err != nil {
		return err
	}
	defer adm.Close()

	resp, err := adm.DeleteTopics(ctx, topic)
	if err != nil {
		return fmt.Errorf("delete topic %q: %w", topic, err)
	}
	for _, r := range resp {
		if r.Err != nil && !isKafkaErr(r.Err, kerr.UnknownTopicOrPartition.Code) {
			return fmt.Errorf("delete topic %q: %w", topic, r.Err)
		}
	}
	return nil
}

// Logger returns a std logger prefixed with the given exercise name so output
// from different exercises is easy to tell apart.
func Logger(exercise string) *log.Logger {
	return log.New(os.Stdout, "["+exercise+"] ", log.LstdFlags|log.Lmsgprefix)
}

// WithTimeout is a small convenience wrapper for the common
// "context with a deadline, plus its cancel" pattern used throughout the labs.
func WithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
