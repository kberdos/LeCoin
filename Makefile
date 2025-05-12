TARGETS = client vswitch
TESTS = test_leminer test_tx test_block
CMD = cmd
MAIN = main.go
TEST_DIR = ./test
TEST_FLAGS = -v -race -count=1

all: $(TARGETS)

test: $(TESTS)

run:
	./util/vnet_run ./configs

client: $(CMD)/client/$(MAIN)
	go build -o $@ $^

vswitch: $(CMD)/vswitch/$(MAIN)
	go build -o $@ $^

test_leminer: $(TEST_DIR)/leminer/*
	go test $(TEST_FLAGS) $(TEST_DIR)/leminer/...

test_tx: $(TEST_DIR)/tx/*
	go test $(TEST_FLAGS) $(TEST_DIR)/tx/...

test_block: $(TEST_DIR)/block/*
	go test $(TEST_FLAGS) $(TEST_DIR)/block/...

clean:
	rm -f client vswitch

.PHONY: all test clean $(TESTS)
