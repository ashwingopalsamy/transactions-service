# transactions-service

A RESTful API service written in Go, utilizing PostgreSQL and Docker to manage accounts and transactions.

[![Go - Build](https://github.com/ashwingopalsamy/transactions-service/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ashwingopalsamy/transactions-service/actions/workflows/go.yml) [![Swagger](https://img.shields.io/badge/Swagger-API%20Spec-white?labelColor=85ea2d&style=plastic&link=https://ashwingopalsamy.github.io/swagger-ui-transactions-service/)](https://ashwingopalsamy.github.io/swagger-ui-transactions-service/) [![Postman](https://img.shields.io/badge/Postman-API%20Doc-white?labelColor=orange&style=plastic&link=https://documenter.getpostman.com/view/19993834/2sAYX8KMKN)](https://documenter.getpostman.com/view/19993834/2sAYX8KMKN)

## Table of Contents

1. [**Architecture and Design**](#1-architecture-and-design)
2. [**Getting Started**](#2-getting-started)
3. [**Endpoints**](#3-endpoints)
4. [**Repository Structure**](#4-repository-structure)
5. [**Running Tests**](#5-running-tests)
6. [**Documentation**](#6-documentation)

## 1. Architecture and Design

The service is structured into **Handler → Service → Repository** layers, ensuring clear separation of concerns and maintainability. It follows a **Domain-Driven Design (DDD)** approach to keep business logic well-structured and scalable.

### **Technical Highlights**
- **PostgreSQL with Goose Migrations** for schema evolution.
- **Docker and Docker Compose** for containerization.
- **Middleware** includes **RequestID tracking, structured logging, and panic recovery** for better observability.
- **Test-Driven Development (TDD)** across layers `testify` for unit testing.


## 2. Getting Started

```sh
# Clone the repository
git clone github.com/ashwingopalsamy/transactions-service
cd transactions-service

# Single-command Script that builds and runs the service
./run.sh

# Else, make use of the Makefile
make run
```

## 3. Endpoints

### Create an Account
```sh
curl -X POST http://localhost:8080/v1/accounts \
     -H "Content-Type: application/json" \
     -d '{"document_number": "12345678900"}'
```
_Response:_
```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

### Retrieve Account Info
```sh
curl -X GET http://localhost:8080/v1/accounts/1
```
_Response:_
```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

### Create a Transaction
```sh
curl -X POST http://localhost:8080/v1/transactions \
     -H "Content-Type: application/json" \
     -d '{"account_id": 1, "operation_type_id": 4, "amount": 123.45}'
```
_Response:_
```json
{
  "transaction_id": 10,
  "account_id": 1,
  "operation_type_id": 4,
  "amount": 123.45,
  "event_date": "2025-02-07T10:32:07Z"
}
```

## 4. Repository Structure

```
transactions-service/
├── cmd/                   # Entrypoint
│   ├── app/               # Main application setup
│   │   ├── main.go        # Application bootstrap
│   │   ├── persistence.go # Database initialization
│   │   ├── server.go      # HTTP server setup
├── internal/              # Core business logic
│   ├── handler/           # API Request Handler Layer
│   │   ├── accounts_handler.go
│   │   ├── transactions_handler.go
│   │   ├── types.go
│   ├── middleware/        # Custom Middlewares
│   │   ├── request_id.go
│   ├── repository/        # Data persistence layer
│   │   ├── accounts_repository.go
│   │   ├── accounts_repository_test.go
│   │   ├── transactions_repository.go
│   │   ├── transactions_repository_test.go
│   │   ├── types.go
│   ├── service/           # Business logic layer
│   │   ├── accounts_service.go
│   │   ├── accounts_service_test.go
│   │   ├── transactions_service.go
│   │   ├── transactions_service_test.go
│   │   ├── types.go
│   ├── writer/            # Response writers
│   │   ├── error_writer.go
├── schema/                # Database schema and migrations
│   ├── migrations/
│   │   ├── 20250207063000_create_trigger_updated_at_timestamp.sql
│   │   ├── 20250207063101_create_table_accounts.sql
│   │   ├── 20250207063202_create_table_operation_types.sql
│   │   ├── 20250207063303_create_table_transactions.sql
│   │   ├── 20250207063404_insert_operation_types_initial_values.sql
│   ├── migrations.Dockerfile
├── docker-compose.yml      # Container orchestration setup
├── Dockerfile              # Service container definition
├── Makefile                # Automation commands for easy dev-experience
├── unittest.Dockerfile     # Containerized unit tests
├── run.sh                  # Single-command script to build and run the service
```

## 5. Running Tests
```sh
make unit-test

(or)

docker-compose run --rm test
```

## 6. Documentation
- SwaggerUI: [OpenAPI v3](https://ashwingopalsamy.github.io/swagger-ui-transactions-service/) Specification
- API Documentation + API Testing : [Postman](https://documenter.getpostman.com/view/19993834/2sAYX8KMKN)


> **Note**: This repository is a **technical demonstration** showcasing my abilities in **API Product Design** and building **scalable, distributed systems**.
>
> It is neither affiliated with; nor endorsed by any organization and does not attempt to replicate any official APIs or services. This project is strictly serves to demonstrate my **technical expertise**.

