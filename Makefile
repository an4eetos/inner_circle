# Inner Circle monorepo tasks. Run `make help` for the list.
SHELL := /bin/bash
.DEFAULT_GOAL := help

-include .env
export

GOTOOL   := go tool -modfile=tools/go.mod
COMPOSE  := docker compose -f infra/docker-compose.yml
OPENAPI_GENERATOR_IMAGE := openapitools/openapi-generator-cli:v7.26.0
FLAVOR   ?= dev

.PHONY: help
help: ## Show this help
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

# ---------------------------------------------------------------- setup

.env:
	cp .env.example .env

.PHONY: setup
setup: .env hooks ## First-time setup: .env, git hooks, dependencies
	cd backend && go mod download && go mod download -modfile=tools/go.mod
	cd web && npm ci
	cd mobile/packages/api_client && dart pub get
	cd mobile && flutter pub get

.PHONY: hooks
hooks: ## Install git hooks (.githooks)
	git config core.hooksPath .githooks

# ---------------------------------------------------------------- local stack

.PHONY: dev
dev: .env infra-up migrate ## Run infra + API + worker with hot reload
	@cd backend && trap 'kill 0' EXIT INT TERM; \
		$(GOTOOL) air -c .air.api.toml & \
		$(GOTOOL) air -c .air.worker.toml & \
		wait

.PHONY: api
api: .env ## Run only the API (hot reload)
	cd backend && $(GOTOOL) air -c .air.api.toml

.PHONY: worker
worker: .env ## Run only the worker (hot reload)
	cd backend && $(GOTOOL) air -c .air.worker.toml

.PHONY: infra-up
infra-up: .env ## Start Postgres, S3 (SeaweedFS), Mailpit
	$(COMPOSE) up -d --wait postgres
	$(COMPOSE) up -d s3 mailpit

.PHONY: infra-down
infra-down: ## Stop the local stack (keeps data)
	$(COMPOSE) down

.PHONY: infra-reset
infra-reset: ## Stop the local stack and DELETE its volumes
	$(COMPOSE) down -v

.PHONY: migrate
migrate: .env ## Apply pending migrations
	cd backend && go run ./cmd/migrate up

.PHONY: migrate-down
migrate-down: .env ## Roll back the latest migration
	cd backend && go run ./cmd/migrate down

.PHONY: migrate-status
migrate-status: .env ## Show migration status
	cd backend && go run ./cmd/migrate status

.PHONY: migration
migration: ## Create a migration: make migration name=add_users
	@test -n "$(name)" || (echo "usage: make migration name=<snake_case>" && exit 1)
	cd backend && $(GOTOOL) goose -dir migrations -s create $(name) sql

.PHONY: web
web: ## Run the Next.js dev server
	cd web && npm run dev

.PHONY: mobile-run
mobile-run: ## Run the Flutter app: make mobile-run FLAVOR=dev|staging|prod
	@# Lets an Android emulator/USB device reach the host's API at localhost:18080.
	-@command -v adb >/dev/null && adb reverse tcp:18080 tcp:18080 >/dev/null 2>&1
	cd mobile && flutter run --flavor $(FLAVOR) --dart-define-from-file=config/$(FLAVOR).json

# ---------------------------------------------------------------- codegen

.PHONY: gen
gen: gen-go gen-dart gen-ts gen-mobile ## Regenerate all code from api/openapi.yaml, SQL, l10n, drift

.PHONY: gen-go
gen-go: ## Go server stubs (oapi-codegen) + typed SQL (sqlc)
	cd backend && $(GOTOOL) oapi-codegen -config oapi-codegen.yaml ../api/openapi.yaml
	cd backend && $(GOTOOL) sqlc generate

.PHONY: gen-dart
gen-dart: ## Dart API client (openapi-generator dart-dio, in Docker)
	rm -rf mobile/packages/api_client/lib mobile/packages/api_client/doc mobile/packages/api_client/test
	docker run --rm -u "$$(id -u):$$(id -g)" -v "$(CURDIR):/local" $(OPENAPI_GENERATOR_IMAGE) generate \
		-i /local/api/openapi.yaml -g dart-dio -o /local/mobile/packages/api_client \
		--additional-properties=pubName=inner_circle_api,pubDescription="Generated Inner Circle API client",serializationLibrary=json_serializable,dateLibrary=core,hideGenerationTimestamp=true,enumUnknownDefaultCase=true \
		--global-property=apiDocs=false,modelDocs=false,apiTests=false,modelTests=false
	cd mobile/packages/api_client && dart pub get && dart run build_runner build --delete-conflicting-outputs
	cd mobile/packages/api_client && dart format lib >/dev/null

.PHONY: gen-ts
gen-ts: ## TypeScript API types (openapi-typescript)
	cd web && npm exec -- openapi-typescript ../api/openapi.yaml -o src/gen/api/schema.d.ts

.PHONY: gen-mobile
gen-mobile: ## Flutter l10n + drift/build_runner outputs
	cd mobile && flutter gen-l10n && dart run build_runner build --delete-conflicting-outputs

# ---------------------------------------------------------------- quality

.PHONY: test
test: test-backend test-mobile test-web ## Run all tests

.PHONY: test-backend
test-backend: ## Go tests (integration tests need Docker)
	cd backend && go test -race ./...

.PHONY: test-mobile
test-mobile: ## Flutter tests
	cd mobile && flutter test

.PHONY: test-web
test-web: ## Web typecheck + build
	cd web && npm run typecheck && npm run build

.PHONY: lint
lint: lint-backend lint-mobile lint-web ## Run all linters

.PHONY: lint-backend
lint-backend:
	cd backend && golangci-lint run ./...

.PHONY: lint-mobile
lint-mobile:
	cd mobile && flutter analyze

.PHONY: lint-web
lint-web:
	cd web && npm run lint

.PHONY: fmt
fmt: ## Format all code
	cd backend && golangci-lint fmt ./...
	cd mobile && dart format lib test
	cd web && npm run format
