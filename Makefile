# Makefile for Go REST API Docker operations

# Variables
APP_NAME := microtester
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "latest")
REGISTRY := docker.io
USERNAME := gigiozzz
IMAGE_NAME := $(REGISTRY)/$(USERNAME)/$(APP_NAME)
TAGGED_IMAGE := $(IMAGE_NAME):$(VERSION)
LATEST_IMAGE := $(IMAGE_NAME):latest

# Docker build arguments
DOCKER_BUILD_ARGS := --platform=linux/amd64,linux/arm64

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: help build build-dev build-prod push push-latest run run-dev stop clean test docker-login check-deps

# Default target
help: ## Show this help message
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Prerequisites check
check-deps: ## Check if required tools are installed
	@echo "$(GREEN)Checking dependencies...$(NC)"
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)Docker is required but not installed$(NC)"; exit 1; }
	@command -v git >/dev/null 2>&1 || { echo "$(RED)Git is required but not installed$(NC)"; exit 1; }
	@echo "$(GREEN)✓ All dependencies found$(NC)"

# Initialize go module if not exists
init: ## Initialize Go module
	@if [ ! -f go.mod ]; then \
		echo "$(YELLOW)Initializing Go module...$(NC)"; \
		go mod init $(APP_NAME); \
	fi
	@go mod tidy

# Local development
test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v ./...

run-dev: ## Run application locally
	@echo "$(GREEN)Starting development server...$(NC)"
	go run main.go

# Docker operations
build: check-deps ## Build Docker image for development
	@echo "$(GREEN)Building Docker image: $(TAGGED_IMAGE)$(NC)"
	docker build -t $(TAGGED_IMAGE) -t $(LATEST_IMAGE) .
	@echo "$(GREEN)✓ Build completed$(NC)"

build-dev: check-deps ## Build Docker image with development settings
	@echo "$(GREEN)Building development Docker image...$(NC)"
	docker build --target builder -t $(IMAGE_NAME):dev .
	@echo "$(GREEN)✓ Development build completed$(NC)"

build-prod: check-deps ## Build optimized production image
	@echo "$(GREEN)Building production Docker image: $(TAGGED_IMAGE)$(NC)"
	docker build \
		--build-arg BUILD_VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		-t $(TAGGED_IMAGE) \
		-t $(LATEST_IMAGE) .
	@echo "$(GREEN)✓ Production build completed$(NC)"

build-multi: check-deps ## Build multi-architecture image
	@echo "$(GREEN)Building multi-architecture image...$(NC)"
	docker buildx create --use --name multiarch-builder 2>/dev/null || true
	docker buildx build \
		$(DOCKER_BUILD_ARGS) \
		-t $(TAGGED_IMAGE) \
		-t $(LATEST_IMAGE) \
		--push .
	@echo "$(GREEN)✓ Multi-architecture build and push completed$(NC)"

# Docker registry operations
push: build docker-login ## Build and push versioned image
	@echo "$(GREEN)Pushing image: $(TAGGED_IMAGE)$(NC)"
	docker push $(TAGGED_IMAGE)
	@echo "$(GREEN)✓ Push completed$(NC)"

push-latest: build-prod ## Push both versioned and latest tags
	@echo "$(GREEN)Pushing latest image: $(LATEST_IMAGE)$(NC)"
	docker push $(LATEST_IMAGE)
	@echo "$(GREEN)✓ Latest push completed$(NC)"

push-all: build-prod docker-login ## Build and push all tags
	@echo "$(GREEN)Pushing all images...$(NC)"
	docker push $(TAGGED_IMAGE)
	docker push $(LATEST_IMAGE)
	@echo "$(GREEN)✓ All images pushed$(NC)"

# Information
info: ## Show build information
	@echo "$(GREEN)Build Information:$(NC)"
	@echo "  App Name: $(APP_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Registry: $(REGISTRY)"
	@echo "  Username: $(USERNAME)"
	@echo "  Tagged Image: $(TAGGED_IMAGE)"
	@echo "  Latest Image: $(LATEST_IMAGE)"
