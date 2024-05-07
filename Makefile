PACKAGE=github.com/tappoy/vault-cli
WORKING_DIRS=tmp bin

SRC=$(shell find . -name "*.go")
BIN=bin/$(shell basename $(CURDIR))
USAGE=Usage.txt
COVER=tmp/cover
COVER0=tmp/cover0

.PHONY: all clean fmt cover test lint

all: $(WORKING_DIRS) $(FMT) $(BIN) test lint

clean:
	rm -rf $(WORKING_DIRS)

$(WORKING_DIRS):
	mkdir -p $(WORKING_DIRS)

fmt: $(SRC)
	go fmt

go.sum: go.mod
	go mod tidy

$(BIN): go.sum $(USAGE)
	go build -o $(BIN)

test: $(BIN)
	go test -v -tags=mock -vet=all -cover -coverprofile=$(COVER)

cover: $(COVER)
	grep "0$$" $(COVER) | sed 's!$(PACKAGE)!.!' | tee $(COVER0)

lint: $(BIN)
	go vet
