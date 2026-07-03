// Exercise 09 — Transactions
//
// Goal: produce records atomically. A transactional producer either makes ALL
// records in a transaction visible (commit) or NONE of them (abort). Consumers
// that read with the ReadCommitted isolation level never see aborted records.
// This is the basis of exactly-once processing.
//
// A transactional producer needs a stable TransactionalID and acks=all.
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

// newTxnProducer returns a client configured for transactions: a stable
// transactional id and acks=all.
//
// TODO: labkafka.NewClient(kgo.TransactionalID(txnID),
// kgo.RequiredAcks(kgo.AllISRAcks())).
func newTxnProducer(txnID string) (*kgo.Client, error) {
	_ = txnID
	return nil, errTODO
}

// produceTxn runs a single transaction: begin, produce every value, then either
// commit (commit == true) or abort (commit == false). On abort, none of the
// records must become visible to ReadCommitted consumers.
//
// TODO:
//  1. cl.BeginTransaction()
//  2. produce all values (ProduceSync, or Produce then Flush)
//  3. cl.EndTransaction(ctx, kgo.TryCommit) or kgo.TryAbort based on `commit`
func produceTxn(ctx context.Context, cl *kgo.Client, topic string, values []string, commit bool) error {
	_ = ctx
	_ = cl
	_ = topic
	_ = values
	_ = commit
	return errTODO
}

func main() {
	log := labkafka.Logger("09-transactions")
	ctx, cancel := labkafka.WithTimeout(30e9)
	defer cancel()

	const topic = "labkafka.09-transactions.demo"
	if err := labkafka.EnsureTopic(ctx, topic, 1, 1); err != nil {
		log.Fatalf("ensure topic: %v", err)
	}

	cl, err := newTxnProducer("labkafka.09-txn.demo")
	if err != nil {
		log.Fatalf("newTxnProducer: %v", err)
	}
	defer cl.Close()

	// This transaction aborts: its records must never be visible.
	if err := produceTxn(ctx, cl, topic, []string{"ghost-1", "ghost-2"}, false); err != nil {
		log.Fatalf("aborted txn: %v", err)
	}
	log.Printf("aborted a transaction (records should be invisible)")

	// This transaction commits: its records become visible.
	if err := produceTxn(ctx, cl, topic, []string{"real-1", "real-2"}, true); err != nil {
		log.Fatalf("committed txn: %v", err)
	}
	log.Printf("committed a transaction")

	fmt.Println("09-transactions: done")
}
