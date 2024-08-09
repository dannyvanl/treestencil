run-example:
    cd example && go run ../...

lint:
    golangci-lint run

build:
    go build -v ./...
