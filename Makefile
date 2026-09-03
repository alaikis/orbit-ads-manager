SHELL := /bin/bash
.PHONY: help dev build test clean migrate up down logs

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start development environment
	docker compose up -d postgres redis
	@echo "PostgreSQL and Redis started. Run 'make dev-backend' and 'make dev-frontend' in separate terminals."

dev-backend: ## Run backend in development
	cd apps/api && go run cmd/server/main.go

dev-frontend: ## Run frontend in development
	cd apps/web && npm install && npm run dev

build: ## Build all services
	cd apps/api && go build ./cmd/server ./cmd/worker ./cmd/scheduler
	cd apps/web && npm run build

test: ## Run tests
	cd apps/api && go test ./...
	cd apps/web && npm run lint && npm run typecheck

migrate: ## Run database migrations
	cd apps/api && go run cmd/server/main.go migrate

up: ## Start all services
	docker compose up -d

down: ## Stop all services
	docker compose down

logs: ## Show logs
	docker compose logs -f

clean: ## Clean build artifacts
	rm -rf apps/api/bin apps/web/.next
