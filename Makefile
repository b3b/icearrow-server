SRC := $(wildcard *.go) $(wildcard */*.go)

all: format test walrus-client icearrow-server

walrus-client: cmd/walrus-client/main.go walrus/walrus.go
	go build -o $@ $<

icearrow-server: cmd/icearrow-server/main.go hippo/hippo.go
	go build -o $@ $<

.PHONY: test
test:
	go test -v ./... -coverprofile=coverage.out

.PHONY: test-e2e
test-e2e:
	go test -v ./... -coverprofile=coverage.out --tags=e2e

.PHONY: format
format:
	go fmt ./...
	go mod tidy -v

.PHONY: clean
clean:
	go clean
	rm -f walrus-client icearrow-server
	rm -f coverage.out
