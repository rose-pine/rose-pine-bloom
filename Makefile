.PHONY: all
all: docs build

.PHONY: build
build:
	go build -o bloom .

.PHONY: docs
docs:
	./script/update-readme-help.sh

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...


.PHONY: check
check:
	golangci-lint run ./...

