WORKING_DIRS=tmp

SRC=$(shell find . -name "*.go")
BIN=tmp/$(shell basename $(CURDIR))
TESTBIN=tmp/$(shell basename $(CURDIR))-test

FMT=tmp/fmt
TEST=tmp/cover

.PHONY: all clean cover

all: $(WORKING_DIRS) $(FMT) $(BIN) $(TEST) $(DOC)

clean:
	rm -rf $(WORKING_DIRS)

$(WORKING_DIRS):
	mkdir -p $(WORKING_DIRS)

$(FMT): $(SRC)
	go fmt ./... > $(FMT) 2>&1 || true

$(BIN): $(SRC)
	go build -o $(BIN)

$(TESTBIN): $(BIN)
	go build -tags test -o $(TESTBIN)

$(TEST): $(TESTBIN)
	go test -v -tags=mock -cover -coverprofile=$(TEST) ./...

cover: $(TEST)
	grep "0$$" $(TEST) || true
