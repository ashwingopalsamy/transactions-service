#!/bin/bash

set -e

echo "Note: Requires Docker and Docker Compose Installed!"

# Variables
SERVICE_NAME="app"
IMPORT_PATH="github.com/ashwingopalsamy/transactions-service"
PKG_SRC="${IMPORT_PATH}/cmd/${SERVICE_NAME}"
BUILD_TAG="local"
IMAGE="${SERVICE_NAME}:${BUILD_TAG}"
MIGRATION_IMAGE="migration_builder:${BUILD_TAG}"
DB_DSN="postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
DOCKER_COMPOSE="docker compose --file docker-compose.yml"

export IMAGE
export MIGRATION_IMAGE

echo "Handling dependencies..."
go mod tidy
go mod vendor

echo "Building the service..."
docker build --pull --tag "${IMAGE}" \
    --build-arg importPath="${IMPORT_PATH}" \
    --build-arg pkg="${PKG_SRC}" .

echo "Building the migration image..."
docker build --tag "${MIGRATION_IMAGE}" \
    --build-arg importPath="${IMPORT_PATH}" \
    -f schema/migrations.Dockerfile .

echo "Running database migrations..."
${DOCKER_COMPOSE} run --rm migrations -dir /migrations postgres "${DB_DSN}" up

echo "Starting all services..."
${DOCKER_COMPOSE} up --build -d
${DOCKER_COMPOSE} logs -f

