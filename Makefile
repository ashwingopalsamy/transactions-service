# Variables
SERVICE_NAME = app
IMPORT_PATH = github.com/ashwingopalsamy/transactions-service
PKG_SRC = ${IMPORT_PATH}/cmd/${SERVICE_NAME}
BUILD_TAG = local
IMAGE = ${SERVICE_NAME}:${BUILD_TAG}
MIGRATION_IMAGE = migration_builder:${BUILD_TAG}
DB_DSN=postgres://postgres:postgres@db:5432/postgres?sslmode=disable
DOCKER_COMPOSE = docker compose --file docker-compose.yml

export IMAGE
export MIGRATION_IMAGE

# Main targets
.PHONY: all build run unit-test migrate clean

all: build run

dep:
	@echo "Handling dependencies..."
	go mod tidy
	go mod vendor

# Build targets
build: build-service build-migrations

build-service:
	@echo "Building the service..."
	docker build --pull --tag ${IMAGE} \
		--build-arg importPath=${IMPORT_PATH} \
		--build-arg pkg=${PKG_SRC} .

build-migrations:
	@echo "Building the migration image..."
	docker build --tag ${MIGRATION_IMAGE} \
		--build-arg importPath=${IMPORT_PATH} \
		-f schema/migrations.Dockerfile .

build-unit-test:
	@echo "Building unit test..."
	docker build --pull --file unittest.Dockerfile \
		--tag ${IMAGE}-unittest \
		--build-arg importPath=${IMPORT_PATH} .

unit-test: build-unit-test
	@echo "Running unit tests..."
	docker run --rm ${IMAGE}-unittest go test -cover -tags=unit ./...


# Docker Compose targets
up:
	@echo "Starting all services..."
	${DOCKER_COMPOSE} up --build -d
	@${DOCKER_COMPOSE} logs -f

down:
	@echo "Stopping all services..."
	${DOCKER_COMPOSE} down

clean:
	@echo "Cleaning up..."
	rm -f ${SERVICE_NAME}
	${DOCKER_COMPOSE} down -v


# Single command: Builds and runs the service with migrations
run: dep build-service build-migrations migrate-up
	@echo "Running the service..."
	${DOCKER_COMPOSE} up --build -d
	@${DOCKER_COMPOSE} logs -f

# Single command: Builds and runs the service without migrations
run-service: dep build-service
	@echo "Running the service..."
	${DOCKER_COMPOSE} up --build -d
	@${DOCKER_COMPOSE} logs -f

