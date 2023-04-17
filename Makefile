.PHONY: all test clean

test:
	go test -coverprofile=coverage.out ./...
cover:
	go tool cover -html=coverage.out