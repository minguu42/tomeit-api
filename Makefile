.PHONY: dev
dev:
	@docker compose up api

.PHONY: down
down:
	@docker compose down

.PHONY: test
test:
	@docker compose up -d db-test
	-go test

.PHONY: fmt
fmt:
	@gofmt -l -s -w .
	@goimports -w .

.PHONY: lint
lint:
	@go vet ./...
	@staticcheck ./...
