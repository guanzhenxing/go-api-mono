.PHONY: all build run test clean check fmt lint vet test-coverage migrate-up migrate-down docker-build docker-run docs swagger integration-test performance-test setup help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=go-api-mono

# Build flags
LDFLAGS=-ldflags "-w -s"

# Environment variables
export GO_ENV ?= development
export CONFIG_PATH ?= configs/config.yaml

all: check build

help:
	@echo "Available commands:"
	@echo "  make build           - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run unit tests"
	@echo "  make check          - Run all checks (fmt, lint, vet, test)"
	@echo "  make clean          - Clean build files"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make vet            - Run go vet"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make migrate-up     - Run database migrations up"
	@echo "  make migrate-down   - Run database migrations down"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run application in Docker"
	@echo "  make docs           - Generate documentation"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make setup          - Initial project setup"
	@echo "  make integration-test - Run integration tests"
	@echo "  make performance-test - Run performance tests"

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/app

run:
	GO_ENV=$(GO_ENV) $(GORUN) ./cmd/app/main.go

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf ./dist
	rm -rf ./coverage

test:
	GO_ENV=testing $(GOTEST) -v ./...

check: fmt lint vet test
	@echo "Running code quality checks..."
	@echo "Checking code formatting..."
	@test -z $$(gofmt -l .)
	@echo "Running linter..."
	golangci-lint run ./...
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo "Running vulnerability check..."
	govulncheck ./...
	@echo "Checking and updating dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "Checking for dependency updates..."
	$(GOCMD) list -u -m all

fmt:
	gofmt -w .
	goimports -w .

lint:
	golangci-lint run

vet:
	$(GOCMD) vet ./...

test-coverage:
	mkdir -p coverage
	GO_ENV=testing $(GOTEST) -coverprofile=coverage/coverage.out ./...
	$(GOCMD) tool cover -html=coverage/coverage.out -o coverage/coverage.html

migrate-up:
	$(GORUN) cmd/migrate/main.go up

migrate-down:
	$(GORUN) cmd/migrate/main.go down

docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker-compose up --build

docs:
	@echo "Generating documentation..."
	mkdir -p docs/api
	swag init -g cmd/app/main.go -o docs/api

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/app/main.go -o api/swagger

setup: check
	@echo "Setting up development environment..."
	$(GOGET) -u github.com/swaggo/swag/cmd/swag
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	@echo "Installing dependencies..."
	$(GOMOD) download
	@echo "Creating necessary directories..."
	mkdir -p logs
	mkdir -p coverage
	@echo "Setup complete!"

integration-test:
	GO_ENV=testing $(GOTEST) -v ./test/integration/...

performance-test:
	GO_ENV=testing $(GORUN) ./test/performance/...

# Development tools installation
tools:
	$(GOGET) -u github.com/swaggo/swag/cmd/swag
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u golang.org/x/vuln/cmd/govulncheck

# Default target
.DEFAULT_GOAL := help