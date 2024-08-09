build:
    go build -v ./...

run-example:
    cd example && go run ../...

lint:
    golangci-lint run
