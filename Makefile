WORKING_DIRS=tmp

SRC=$(shell find . -name "*.go")
BIN=tmp/$(shell basename $(CURDIR))
USAGE=Usage.txt
COVER=tmp/cover
COVER0=tmp/cover0

.PHONY: all clean fmt cover test lint testlint

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

$(COVER0): $(COVER)
	grep "0$$" $(COVER) | tee > $(COVER0) 2>&1

cover: $(COVER)
	go tool cover -html=$(COVER)

lint: $(BIN)
	go vet
