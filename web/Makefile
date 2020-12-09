# Basic go commands
GOCMD=go
GOGET=$(GOCMD) get
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

all: gofmt linter

install:
	go build -o $(GOPATH)/bin/web-glinv

gofmt:
	gofmt -w .

test:
	go test -count=1 -race -cover -v ./...

linter:
	golangci-lint run --enable-all

build:
	go build -o web-glinv

clean:
	go clean -i ./...