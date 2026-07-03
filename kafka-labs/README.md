# kafka-labs

A hands-on, test-driven tour of Apache Kafka in Go using
[franz-go](https://github.com/twmb/franz-go). Each exercise isolates **one**
concept. The shared plumbing (client/admin factory, topic management, logging)
is done for you in [`internal/labkafka`](internal/labkafka/labkafka.go); your job
is to fill in the `TODO` sections in each exercise's `main.go`.

## How the exercises work

Every exercise directory contains:

- **`main.go`** — a runnable program with helper wiring already in place and the
  core Kafka calls left as `TODO`s. Stubs return a "not implemented" error so the
  package still **compiles** out of the box.
- **`main_test.go`** — an integration test that acts as the **spec** for the
  exercise. It talks to a live broker, so it *skips* when `KAFKA_BROKERS` is
  unset and *fails* (red) until you complete the TODOs.

The intended loop:

```
cd exercises/01-admin
go test ./...        # red — see what behaviour is expected
# edit main.go, fill in the TODOs
go test ./...        # green — you implemented the concept
go run .             # watch it work end to end
```

## Prerequisites

- Go 1.25+
- A reachable Kafka broker. Point the lab at it with:

  ```sh
  export KAFKA_BROKERS=localhost:9092   # comma-separated for a cluster
  ```

  If unset, the code defaults to `localhost:9092`, but **tests skip** unless the
  variable is explicitly set (so `go test ./...` is green on a machine with no
  Kafka).

## The progression

| # | Exercise | Concept |
|---|----------|---------|
| 01 | `01-admin` | Create / list / describe / delete topics with `kadm` |
| 02 | `02-produce` | Produce records; sync vs async and the produce callback |
| 03 | `03-consume` | Consume without a group; `PollFetches` and iterating records |
| 04 | `04-consumer-groups` | Join a consumer group; auto vs manual commit |
| 05 | `05-offsets` | Offset management: seek, commit specific offsets, read committed |
| 06 | `06-partitioning` | Keys, partitions, and record ordering guarantees |
| 07 | `07-serialization` | Encode/decode structured values (JSON) on the wire |
| 08 | `08-idempotent-delivery` | The idempotent producer — no duplicates on retry |
| 09 | `09-transactions` | Transactional produce/consume for exactly-once |
| 10 | `10-dlq-retry` | Dead-letter queue + retry topic pattern for poison records |

## Handy Make targets

```sh
make test              # go test ./... across the whole lab
make run EX=02-produce # go run a single exercise
make topics            # list topics on the broker
make tidy              # go mod tidy
```

See the [`Makefile`](Makefile) for the full list.
