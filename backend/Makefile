BUILD_TS := $(shell date -Iseconds -u | sed 's/+00:00/Z/')
COMMIT_SHA ?= $(shell git rev-parse HEAD)
VERSION ?= $(shell git describe --abbrev=0 --tags || echo "latest")

export CGO_ENABLED=0

binary := checkin-system

module=$(shell go list -m)
ld_flags := "-X $(module)/pkg/version.Version=$(VERSION) -X $(module)/pkg/version.GitCommit=$(COMMIT_SHA) -X $(module)/pkg/version.BuildTime=$(BUILD_TS)"

FILES    := $(shell find . -name '*.go' -type f -not -name '*.pb.go' -not -name '*_generated.go' -not -name '*_test.go')
TESTS    := $(shell find . -name '*.go' -type f -not -name '*.pb.go' -not -name '*_generated.go' -name '*_test.go')

.DEFAULT_GOAL := all
.PHONY: all
all: generate fmt lint build

.PHONY: generate
generate:
	go generate ./...

.PHONY: fmt
fmt:
	golangci-lint fmt

.PHONY: lint
lint:
	golangci-lint run

.PHONY: it
it:
	gotestsum -- -v -tags=integration ./...

.PHONY: build
build:
	go build -ldflags $(ld_flags) -o $(binary) cmd/backend/main.go

.PHONY: clean
clean:
	rm -f $(binary)
	go clean -testcache


