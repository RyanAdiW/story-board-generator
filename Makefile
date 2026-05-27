ifneq (,$(wildcard ./.env))
include .env
export
endif

GO ?= go

.PHONY: check-go run run-api run-worker build build-api build-worker tidy

check-go:
	@command -v $(GO) >/dev/null 2>&1 || (echo "go not found in PATH. Install Go in WSL or run with GO=<binary>"; exit 1)

run: run-api

run-api: check-go
	$(GO) run ./cmd/api

run-worker: check-go
	$(GO) run ./cmd/worker

build: build-api build-worker

build-api: check-go
	$(GO) build ./cmd/api

build-worker: check-go
	$(GO) build ./cmd/worker

tidy: check-go
	$(GO) mod tidy
