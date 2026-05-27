ifneq (,$(wildcard ./.env))
include .env
export
endif

GO ?= go

.PHONY: check-go run build tidy

check-go:
	@command -v $(GO) >/dev/null 2>&1 || (echo "go not found in PATH. Install Go in WSL or run with GO=<binary>"; exit 1)

run: check-go
	$(GO) run ./cmd/api

build: check-go
	$(GO) build ./cmd/api

tidy: check-go
	$(GO) mod tidy
