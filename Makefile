.PHONY: test lint ci

test:
	go test ./...

lint:
	go vet ./...

ci: lint test
make ci

