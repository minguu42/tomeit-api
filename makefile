.PHONY:test
test:
	@docker compose up -d db-test
	@sleep 15
	-go test
	@docker compose down db-test
	@docker volume rm tomeit-api_data-test

.PHONY: fmt
fmt:
	@gofmt -l -s -w .
	@goimports -w .

.PHONY: lint
lint:
	@go vet ./...
