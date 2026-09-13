# xermess — authentication server
#
# Run `make` to see every target.
#
# First time here?  make setup

# ---------------------------------------------------------------------------
# Settings
# ---------------------------------------------------------------------------

BINARY     := xermess
BIN_DIR    := bin
SERVER     := ./cmd/xermess
MIGRATE    := ./cmd/migrate
MIGRATIONS := migrations
WEB_DIR    := web

# Go packages, named explicitly: ./... would walk into web/node_modules, which
# contains a stray Go package.
PKGS := ./cmd/... ./internal/... ./migrations/...

# Read .env so the database targets know which database to talk to. These stay
# make variables and are deliberately not exported: exporting would put the
# database password into the environment of every command make runs, including
# the tests, which then would not be testing what they think they are.
ifneq (,$(wildcard .env))
include .env
endif

DB_URL ?= $(XERMESS_DB_DSN)

# The migrate command reads .env itself. Passing the database in keeps it on
# the same database as the psql targets below when DB_URL is overridden on the
# command line, and changes nothing when it is not.
MIGRATE_DB := $(if $(DB_URL),XERMESS_DB_DSN="$(DB_URL)")

# The database name, and the same server with the "postgres" database instead:
# creating a database means connecting to a different one.
DB_NAME  := $(shell echo "$(DB_URL)" | sed -E 's|.*/([^/?]+)(\?.*)?$$|\1|')
DB_ADMIN := $(shell echo "$(DB_URL)" | sed -E 's|/[^/?]+(\?.*)?$$|/postgres\1|')

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "; print "\nUsage: make <target>\n"} \
		/^# -+$$/ {next} \
		/^## / {printf "\n\033[1m%s\033[0m\n", substr($$0, 4); next} \
		/^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' \
		$(MAKEFILE_LIST)
	@echo ""

## Getting started

.PHONY: setup
setup: env deps db-create migrate-up ## Set the project up from scratch, then say what to do next
	@echo ""
	@echo "Ready. Start the server with:  make run"
	@echo "The frontend needs one more step:  make web-install"

.PHONY: env
env: ## Create .env from .env.example if it is missing
	@test -f .env && echo ".env already exists, leaving it alone" \
		|| { cp .env.example .env; echo "created .env — check the database settings in it"; }

.PHONY: deps
deps: ## Download Go dependencies
	go mod download

## Development

.PHONY: run
run: ## Run the server
	go run $(SERVER)

.PHONY: run-migrate
run-migrate: ## Run the server, applying migrations first
	XERMESS_DB_MIGRATE=true go run $(SERVER)

.PHONY: build
build: ## Build the server into bin/
	go build -o $(BIN_DIR)/$(BINARY) $(SERVER)

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
	rm -rf $(WEB_DIR)/build $(WEB_DIR)/.svelte-kit
	go clean

## Database
#  Migrations are Go files in migrations/. They are Go functions, so only a
#  binary that imports them can run them: that is ./cmd/migrate, not the goose
#  command-line tool.

.PHONY: migrate-up
migrate-up: ## Apply pending migrations
	$(MIGRATE_DB) go run $(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## Roll the newest migration back
	$(MIGRATE_DB) go run $(MIGRATE) down

.PHONY: migrate-status
migrate-status: ## Show which migrations are applied
	$(MIGRATE_DB) go run $(MIGRATE) status

.PHONY: migrate-new
migrate-new: ## Create an empty migration — make migrate-new name=add_x
	@test -n "$(name)" || { echo "usage: make migrate-new name=add_something"; exit 1; }
	go tool goose -dir $(MIGRATIONS) create $(name) go

.PHONY: db-create
db-create: ## Create the database named in .env, if it does not exist
	@psql "$(DB_ADMIN)" -tAc "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" | grep -q 1 \
		&& echo "database $(DB_NAME) already exists" \
		|| { psql "$(DB_ADMIN)" -c 'CREATE DATABASE "$(DB_NAME)"' >/dev/null && echo "created database $(DB_NAME)"; }

.PHONY: db-reset
db-reset: ## Delete everything in the database and migrate from scratch
	@printf 'This deletes every table and row in "%s". Type yes to continue: ' "$(DB_NAME)"
	@read answer && [ "$$answer" = "yes" ] || { echo "cancelled"; exit 1; }
	psql "$(DB_URL)" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	$(MAKE) migrate-up

## Quality

.PHONY: check
check: fmt vet test ## Format, vet and test — run this before pushing

.PHONY: fmt
fmt: ## Format Go code
	go fmt $(PKGS)

.PHONY: vet
vet: ## Report suspicious code
	go vet $(PKGS)

.PHONY: test
test: ## Run tests
	go test $(PKGS)

.PHONY: tidy
tidy: ## Add missing and remove unused modules
	go mod tidy

## Frontend

.PHONY: web-install
web-install: ## Install frontend dependencies
	cd $(WEB_DIR) && bun install

.PHONY: web-dev
web-dev: ## Run the frontend dev server
	cd $(WEB_DIR) && bun run dev

.PHONY: web-build
web-build: ## Build the frontend
	cd $(WEB_DIR) && bun run build

.PHONY: web-start
web-start: web-build ## Build the frontend and serve it on :4173
	cd $(WEB_DIR) && bun run preview
