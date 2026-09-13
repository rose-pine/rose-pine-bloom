.PHONY: all
all: build

.PHONY: build
build:
	go build

.PHONY: install
install:
	go install

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	./scripts/golangci-lint-shim.sh run
