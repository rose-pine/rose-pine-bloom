.PHONY: all
all: format build test

.PHONY: docs
docs:
	go generate

.PHONY: build
build:
	go build -o bloom .

.PHONY: test
test:
	go test ./...

.PHONY: format
format:
	go fmt ./...

.PHONY: check
check:
	golangci-lint run ./...

.PHONY: clean
clean:
	rm -f bloom
	rm -rf docs/*.md
