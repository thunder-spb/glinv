.PHONY: build
build:
	go build -v ./cmd/api

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: linter
linter:
	golangci-lint run --enable-all

.DEFAULT_GOAL := build