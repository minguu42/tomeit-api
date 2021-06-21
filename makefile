.PHONY: fmt
fmt:
	@gofmt -l -s -w .
	@goimports -w .

.PHONY: lint
lint:
	go vet ./...
