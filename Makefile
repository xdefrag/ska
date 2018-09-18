default: test test-benchmark test-race

build:
	go build -o ./dist/ska ./cmd/ska

test:
	go test -coverprofile=coverage.txt -covermode=atomic ./...

test-benchmark:
	go test -bench=. ./...

test-race:
	go test -race ./...
