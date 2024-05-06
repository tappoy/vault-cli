WORKING_DIRS=tmp

SRC=$(shell find . -name "*.go")
BIN=tmp/$(shell basename $(CURDIR))
USAGE=Usage.txt
COVER=tmp/cover
COVER0=tmp/cover0

.PHONY: all clean fmt cover test lint testlint

all: $(WORKING_DIRS) $(FMT) $(LINT) $(BIN) $(TEST)

clean:
	rm -rf $(WORKING_DIRS)

$(WORKING_DIRS):
	mkdir -p $(WORKING_DIRS)

fmt: $(SRC)
	go fmt

lint: $(SRC)
	go vet

go.sum: go.mod
	go mod tidy

$(BIN): lint go.sum $(USAGE)
	go build -o $(BIN)

test: $(BIN) $(COVER) testlint
	go test -v -tags=mock -vet=all -cover -coverprofile=$(COVER)

$(COVER0): $(COVER)
	grep "0$$" $(COVER) | tee > $(COVER0) 2>&1

cover: $(COVER)
	go tool cover -html=$(COVER)

