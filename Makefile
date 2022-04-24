build:
	go build

test:
	go test ./... -count=1

lint:
	golangci-lint run -c .golangci.yml