
# Variables
APP_NAME := go-gin-api

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Build the application
build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

# Create a new migration file
# Usage: make migrate-create name=create_users_table
migrate-create:
	@echo "Creating migration files for ${name}..."
	@mkdir -p migrations
	@bash -c 'timestamp=$(date +%Y%m%d%H%M%S); \
	touch migrations/$${timestamp}_${name}.up.sql; \
	touch migrations/$${timestamp}_${name}.down.sql; \
	echo "Created migrations/$${timestamp}_${name}.up.sql"; \
	echo "Created migrations/$${timestamp}_${name}.down.sql"'

# Database migrations
# Load env variables if .env exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

MIGRATE_URL := "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)"

# Apply all up migrations
migrate-up:
	migrate -database $(MIGRATE_URL) -path migrations up

# Apply all down migrations
migrate-down:
	@if [ "$(APP_ENV)" = "production" ] && [ "$(CONFIRM)" != "yes" ]; then \
		echo "ERROR: DANGER! You are attempting to rollback in PRODUCTION."; \
		echo "This can cause DATA LOSS. To proceed, run: make migrate-down CONFIRM=yes"; \
		exit 1; \
	fi
	migrate -database $(MIGRATE_URL) -path migrations down

# Force migration version
# Usage: make migrate-force version=1
migrate-force:
	migrate -database $(MIGRATE_URL) -path migrations force $(version)

.PHONY: run test build migrate-create migrate-up migrate-down migrate-force
