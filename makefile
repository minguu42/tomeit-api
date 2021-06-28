.PHONY:dev
dev:
	@docker compose up -d

.PHONY:stop
stop:
	@docker compose down

.PHONY:log
log:
	@docker compose logs --tail 20

.PHONY:test
test:
	@docker compose up -d db-test
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
