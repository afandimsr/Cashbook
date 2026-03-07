#!/bin/bash

# ==============================================================================
# CashBook Deployment Script
# ==============================================================================
# This script automates the deployment of the CashBook application using
# Docker Compose. It handles network setup, building images, and migrations.
# ==============================================================================

# Configuration - Primary .env file path
ROOT_ENV="./.env"

# Load root .env if it exists
if [ -f "$ROOT_ENV" ]; then
    export $(grep -v '^#' "$ROOT_ENV" | xargs)
fi

# Configuration with defaults
NETWORK_NAME="${NETWORK_NAME:-bridge}"
FRONTEND_ENV="${FRONTEND_ENV_PATH:-./frontend/.env}"
BACKEND_ENV="${BACKEND_ENV_PATH:-./backend/.env}"
STABILIZE_WAIT="${STABILIZE_WAIT:-10}"
RUN_MIGRATIONS="${RUN_MIGRATIONS:-true}"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting CashBook Deployment...${NC}"

# 1. Check for required environment files
echo -e "${YELLOW}🔍 Checking environment files...${NC}"
MISSING_ENV=false
for env_file in "$ROOT_ENV" "$BACKEND_ENV" "$FRONTEND_ENV"; do
    if [ ! -f "$env_file" ]; then
        echo -e "${RED}❌ Error: Missing $env_file${NC}"
        MISSING_ENV=true
    fi
done

if [ "$MISSING_ENV" = true ]; then
    echo -e "${RED}Please create the missing .env files before deploying.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Environment files found.${NC}"

# 2. Ensure the external network exists
echo -e "${YELLOW}🌐 Checking Docker network: $NETWORK_NAME...${NC}"
if ! docker network inspect "$NETWORK_NAME" >/dev/null 2>&1; then
    echo -e "${YELLOW}Creating external network: $NETWORK_NAME...${NC}"
    docker network create "$NETWORK_NAME"
else
    echo -e "${GREEN}✅ Network $NETWORK_NAME already exists.${NC}"
fi

# 3. Stop existing containers (optional, ensures clean state)
# echo -e "${YELLOW}🛑 Stopping existing components...${NC}"
# docker compose down

# 4. Build and start the containers
echo -e "${YELLOW}🛠️  Building and starting containers...${NC}"
# We include both the root .env (for migration vars) and frontend .env (for build args)
# Note: Multiple --env-file flags require Docker Compose V2. If older, consider merging envs.
if docker compose version | grep -q "v2"; then
    docker compose --env-file "$FRONTEND_ENV" up -d --build frontend
else
    # Fallback for older versions: use root .env and assume frontend vars might be needed manually 
    # or just use the one that covers build args.
    echo -e "${YELLOW}⚠️  Note: Docker Compose V1 detected. Using frontend .env for build args.${NC}"
    docker compose --env-file "$FRONTEND_ENV" up -d --build frontend
fi

# 5. Wait for database to be ready
echo -e "${YELLOW}⏳ Waiting for services to stabilize (${STABILIZE_WAIT}s)...${NC}"
sleep "$STABILIZE_WAIT"

# 6. Run database migrations
if [ "$RUN_MIGRATIONS" = "true" ]; then
    echo -e "${YELLOW}🐘 Running database migrations...${NC}"
    if docker compose --env-file "$ROOT_ENV" run --rm migrate; then
        echo -e "${GREEN}✅ Migrations completed successfully.${NC}"
    else
        echo -e "${RED}❌ Error: Migrations failed!${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⏩ Skipping database migrations (RUN_MIGRATIONS != true).${NC}"
fi

echo -e "\n${GREEN}✨ Deployment Successful! ✨${NC}"
echo -e "Backend: http://localhost:8181"
echo -e "Frontend: http://localhost:3000"
echo -e "=============================================================================="
