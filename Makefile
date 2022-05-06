.PHONY: wire build api

api:
	go run github.com/bufbuild/buf/cmd/buf@latest generate

wire:
	cd cmd/example/wired; go run github.com/google/wire/cmd/wire@latest

# Build binary
HEAD = $(shell git rev-parse HEAD)
BUILD_DATE = $(shell date +%FT%T%z)
CI_PIPELINE_ID ?= unknown
LDFLAGS = -X main.version=${HEAD} -X main.buildDate=${BUILD_DATE} -X main.buildNumber=1.0.${CI_PIPELINE_ID}
GOOS=linux
GOARCH=amd64

build: wire
	@echo "Build binary..."
	go build -ldflags "$(LDFLAGS)" -o ./bin/example ./cmd/example/main.go
