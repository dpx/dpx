GOOS=darwin
GOARCH=amd64

.PHONY: build
build:
	dpx exec -e GOOS=$(GOOS) -e GOARCH=$(GOARCH) go build -o ./bin/ ./cmd/...

.PHONY: test
test:
	dpx exec go test ./...

.PHONY: clean
clean:
	rm -rf bin/