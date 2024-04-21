SRC=$(shell find . -name "*.go")
BIN=bin/$(shell basename $(CURDIR))
TESTBIN=bin/$(shell basename $(CURDIR))-test

all: fmt test build

fmt:
	go fmt ./...

build: $(SRC)
	go build -o $(BIN)

build-test: $(SRC)
	go build -tags test -o $(TESTBIN)

test: build-test
	go test -v ./...
