default: test test-benchmark test-race

build:
	go build -o ./dist/ska ./cmd/ska

test:
	go test -count=1 -covermode=count ./...

test-benchmark:
	go test -count=1 -bench=. ./...

test-race:
	go test -count=1 -race ./...
