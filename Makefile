.PHONY: all
all: format build test

.PHONY: build
build:
	go build -o bloom .

.PHONY: test
test:
	go test ./...

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -f bloom
