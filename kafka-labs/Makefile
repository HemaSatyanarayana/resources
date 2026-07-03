# kafka-labs — common tasks
#
# Point the lab at your broker first (defaults to localhost:9092 in code, but
# tests only run when KAFKA_BROKERS is set):
#
#   export KAFKA_BROKERS=localhost:9092

KAFKA_BROKERS ?= localhost:9092
export KAFKA_BROKERS

# Directory name of a single exercise, e.g. EX=02-produce
EX ?= 01-admin

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Compile every exercise and the helper package
	go build ./...

.PHONY: vet
vet: ## go vet the whole module
	go vet ./...

.PHONY: test
test: ## Run all tests (integration tests hit KAFKA_BROKERS)
	go test ./...

.PHONY: test-unit
test-unit: ## Run tests with no broker (integration tests skip)
	env -u KAFKA_BROKERS go test ./...

.PHONY: run
run: ## Run one exercise: make run EX=02-produce
	go run ./exercises/$(EX)

.PHONY: topics
topics: ## List topics on the broker (needs the rpk or kafka CLI; falls back to a Go lister)
	@go run ./exercises/01-admin >/dev/null 2>&1 || true
	@echo "Use exercise 01-admin's kadm calls, or your broker's CLI, to inspect topics."

.PHONY: tidy
tidy: ## go mod tidy
	go mod tidy

.PHONY: fmt
fmt: ## gofmt the tree
	gofmt -w .
