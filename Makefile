# Default target: run Go tests
.PHONY: all
all: test

# Run all Go tests
.PHONY: test
test:
	go test ./...

# Format and tidy the code
.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w .
	go mod tidy

# Clean generated test files
.PHONY: clean
clean:
	rm -rf testdata/*.tfrecord
