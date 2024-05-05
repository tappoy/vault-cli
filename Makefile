WORKING_DIRS=tmp

SRC=$(shell find . -name "*.go")
BIN=tmp/$(shell basename $(CURDIR))
TESTBIN=tmp/$(shell basename $(CURDIR))-test
USAGE=Usage.txt

FMT=tmp/fmt
TEST=tmp/cover

.PHONY: all clean cover test

all: $(WORKING_DIRS) $(FMT) $(BIN) $(TEST)

clean:
	rm -rf $(WORKING_DIRS)

$(WORKING_DIRS):
	mkdir -p $(WORKING_DIRS)

$(FMT): $(SRC)
	go fmt ./... > $(FMT) 2>&1 || true

go.sum: go.mod
	go mod tidy

$(BIN): $(SRC) go.sum $(USAGE)
	go build -o $(BIN)

$(TESTBIN): $(BIN)
	go build -tags test -o $(TESTBIN)

$(TEST): $(TESTBIN)
	make test

test:
	go test -v -tags=test -cover -coverprofile=$(TEST) ./...

cover: $(TEST)
	grep "0$$" $(TEST) || true
