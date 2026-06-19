.PHONY: check test

check:
	golangci-lint run ./...

test:
	go test ./...
