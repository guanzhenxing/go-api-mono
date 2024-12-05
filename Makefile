.PHONY: all build clean test coverage lint run migrate-up migrate-down docker-build docker-run help check secure update-tools check-go-version

# 变量定义
APP_NAME := go-api-mono
GO := go
GOFLAGS := -v
BINARY := $(APP_NAME)
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"
COVERAGE_DIR := coverage
GOLANGCI_LINT_VERSION := v1.55.2

# 默认环境
ENV ?= development

# Docker相关变量
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG ?= latest

# 帮助信息
help:
	@echo "Management commands for $(APP_NAME):"
	@echo
	@echo "Usage:"
	@echo "    make build          Build binary"
	@echo "    make run            Run application"
	@echo "    make clean          Clean build files"
	@echo "    make test           Run tests"
	@echo "    make coverage       Run tests with coverage"
	@echo "    make lint           Run linter"
	@echo "    make check          Run all checks (format, lint, test, vet)"
	@echo "    make bench          Run benchmarks"
	@echo "    make vuln           Run vulnerability check"
	@echo "    make migrate-up     Run database migrations up"
	@echo "    make migrate-down   Run database migrations down"
	@echo "    make docker-build   Build docker image"
	@echo "    make docker-run     Run docker container"
	@echo "    make swagger        Generate swagger documentation"
	@echo "    make help           Show this help message"
	@echo
	@echo "Environment variables:"
	@echo "    ENV                 Set environment (development|production|testing)"
	@echo "    DOCKER_TAG         Set docker tag (default: latest)"

# 构建
build:
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY) ./cmd/app

# 运行
run:
	GO_ENV=$(ENV) $(GO) run ./cmd/app/main.go

# 清理
clean:
	rm -rf bin/
	rm -rf $(COVERAGE_DIR)/
	find . -type f -name '*.out' -delete
	find . -type f -name '*.test' -delete
	find . -type f -name '*.prof' -delete
	find . -type f -name '*.cov' -delete

# 测试
test:
	GO_ENV=testing $(GO) test -v ./...

# 测试覆盖率
coverage:
	@mkdir -p $(COVERAGE_DIR)
	GO_ENV=testing $(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated in $(COVERAGE_DIR)/coverage.html"
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out | grep total | awk '{print "Total coverage: " $$3}'

# 代码检查
check: fmt lint vet test staticcheck secure

# 格式化检查
fmt:
	@echo "Checking code formatting..."
	@test -z $$($(GO) fmt ./...)

# 代码检查
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	fi
	golangci-lint run ./...

# 静态分析
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# 静态检查工具
staticcheck:
	@echo "Running staticcheck..."
	@if ! command -v staticcheck &> /dev/null; then \
		echo "Installing staticcheck..." && \
		$(GO) install honnef.co/go/tools/cmd/staticcheck@latest; \
	fi
	staticcheck ./...

# 基准测试
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# 性能分析
profile:
	@mkdir -p $(COVERAGE_DIR)
	$(GO) test -cpuprofile=$(COVERAGE_DIR)/cpu.prof -memprofile=$(COVERAGE_DIR)/mem.prof -bench=. ./...
	$(GO) tool pprof -http=:2024 $(COVERAGE_DIR)/cpu.prof

# 漏洞检查
vuln:
	@echo "Running vulnerability check..."
	@if ! command -v govulncheck &> /dev/null; then \
		echo "Installing govulncheck..." && \
		$(GO) install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
	govulncheck ./...
	@echo "Vulnerability check completed."

# 依赖检查和更新
deps:
	@echo "Checking and updating dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	$(GO) mod verify
	@echo "Checking for dependency updates..."
	$(GO) list -u -m all
	@echo "Dependencies check completed."

# 安全检查（包含所有安全相关的检查）
secure: vuln deps lint staticcheck
	@echo "All security checks completed."

# 更新所有工具到最新版本
update-tools:
	@echo "Updating development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@echo "Development tools updated."

# 检查 Go 版本
check-go-version:
	@echo "Checking Go version..."
	@go version
	@echo "Required Go version: $(shell grep -E '^go [0-9]+\.[0-9]+(\.[0-9]+)?$$' go.mod | cut -d' ' -f2)"

# 数据库迁移
migrate-up:
	$(GO) run cmd/migrate/main.go up

migrate-down:
	$(GO) run cmd/migrate/main.go down

# Docker构建
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker运行
docker-run:
	docker run -p 2024:2024 --env-file .env.$(ENV) $(DOCKER_IMAGE):$(DOCKER_TAG)

# 生成swagger文档
swagger:
	@if ! command -v swag &> /dev/null; then \
		echo "Installing swag..." && \
		$(GO) install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init -g cmd/app/main.go -o api/swagger

# 安装工具
tools:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest

# 默认目标
all: check build