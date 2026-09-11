BINARY := xermess
BIN_DIR := bin
CMD := ./cmd/xermess
WEB_DIR := web

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the server into bin/
	go build -o $(BIN_DIR)/$(BINARY) $(CMD)

.PHONY: run
run: ## Run the server
	go run $(CMD)

.PHONY: test
test: ## Run Go tests
	go test ./...

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: tidy
tidy: ## Tidy go.mod
	go mod tidy

.PHONY: web-install
web-install: ## Install frontend dependencies
	cd $(WEB_DIR) && bun install

.PHONY: web-dev
web-dev: ## Run the frontend dev server
	cd $(WEB_DIR) && bun run dev

.PHONY: web-build
web-build: ## Build the frontend
	cd $(WEB_DIR) && bun run build

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
	rm -rf $(WEB_DIR)/build $(WEB_DIR)/.svelte-kit
	go clean
