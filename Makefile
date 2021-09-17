.DEFAULT_GOAL := help

.PHONY: dev
dev: ## ローカルサーバを起動する
	@docker compose up api

.PHONY: down
down: ## 関連する Docker コンテナを削除する
	@docker compose down

.PHONY: docs
docs: ## Swagger UI で http://localhost:8000 に仕様書を表示する
	@docker compose up -d docs

.PHONY: fmt
fmt: ## gofmt, goimports による自動整形を実行する
	@gofmt -l -s -w .
	@goimports -w .

.PHONY: lint
lint: ## govet, staticcheck による静的解析を実行する
	@go vet ./...
	@staticcheck ./...

.PHONY: test
test:  ## テストを実行する
	@docker compose up -d db-test
	@go test

.PHONY: check
check: ## fmt, lint, test を実行し, 適切な状態か確認する
	@make fmt
	@make lint
	@make test

.PHONY: bench
bench: ## ハンドラ関数のベンチマークを実行する
	@go test -bench .

.PHONY: help
help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'
