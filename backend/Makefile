
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
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	touch migrations/$$timestamp_${name}.up.sql; \
	touch migrations/$$timestamp_${name}.down.sql; \
	echo "Created migrations/$$timestamp_${name}.up.sql"; \
	echo "Created migrations/$$timestamp_${name}.down.sql"

# Database migrations
# Load env variables if .env exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Database URL construction (supports postgres and mysql, mirrors migrate.ps1)
ifneq ($(DATABASE_URL),)
    MIGRATE_URL := $(DATABASE_URL)
else
    DB_DRIVER ?= postgres
    ifeq ($(DB_DRIVER),mysql)
        MIGRATE_URL := "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)"
    else
        MIGRATE_URL := "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"
    endif
endif

# Apply all up migrations
# Usage: make migrate-up [n=1]
migrate-up:
	migrate -database $(MIGRATE_URL) -path migrations up $(n)

# Apply all down migrations
# Usage: make migrate-down [n=1]
migrate-down:
	@if [ "$(APP_ENV)" = "production" ] && [ "$(CONFIRM)" != "yes" ]; then \
		echo "ERROR: DANGER! You are attempting to rollback in PRODUCTION."; \
		echo "This can cause DATA LOSS. To proceed, run: make migrate-down CONFIRM=yes"; \
		exit 1; \
	fi
	migrate -database $(MIGRATE_URL) -path migrations down $(n)

# Force migration version
# Usage: make migrate-force version=1
migrate-force:
	migrate -database $(MIGRATE_URL) -path migrations force $(version)

# Generate swagger documentation
swag:
	swag init -g cmd/api/main.go -o docs

.PHONY: run test build migrate-create migrate-up migrate-down migrate-force swag
