build:
    go build ./...

run-example:
    cd example && go run ../...

errcheck:
    go install github.com/kisielk/errcheck@latest
    ~/go/bin/errcheck ./...

