.PHONY:dev
dev:
	@docker compose --file ./build/docker-compose.yaml up -d

.PHONY:stop
stop:
	@docker compose --file ./build/docker-compose.yaml down

.PHONY:log
log:
	@docker compose --file ./build/docker-compose.yaml logs --tail 20

.PHONY: fmt
fmt:
	@gofmt -l -s -w .
	@goimports -w .

.PHONY: lint
lint:
	go vet ./...
