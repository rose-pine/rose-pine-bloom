.Phony: all
all: docs build

.Phony: build
build:
	go build -o bloom .

.Phony: docs
docs:
	./script/update-readme-help.sh

.Phony: test
test:
	go test ./...

