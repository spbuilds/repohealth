BINARY=repohealth
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/spbuilds/repohealth/internal/cli.Version=$(VERSION)"

.PHONY: build test lint install clean

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/repohealth/

test:
	go test ./... -v

lint:
	golangci-lint run ./...

install:
	go install $(LDFLAGS) ./cmd/repohealth/

clean:
	rm -f $(BINARY)
	go clean
