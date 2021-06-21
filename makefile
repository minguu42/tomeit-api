.PHONY: fmt
fmt:
	gofmt -l -s -w ./...

.PHONY: lint
lint:
	go vet ./...
