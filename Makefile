# Variables
APP_NAME := embed-migrations
MAIN_PATH := ./main.go
BINARY_PATH := ./bin/$(APP_NAME)

# Go related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)

# Docker related variables
DOCKER_IMAGE := your-docker-image-name
DOCKER_TAG := latest

# Make sure to create bin directory
$(shell mkdir -p $(GOBIN))

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Display available commands
	@echo "Available commands:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the application
	go build -o $(BINARY_PATH) $(MAIN_PATH)

.PHONY: run
run: ## Run the application
	go run $(MAIN_PATH)

.PHONY: clean
clean: ## Clean build files
	rm -rf $(GOBIN)
	go clean

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -cover ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run: ## Run docker container
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: migrate-up
migrate-up: ## Run database migrations up
	go run $(MAIN_PATH) migrate up

.PHONY: migrate-down
migrate-down: ## Run database migrations down
	go run $(MAIN_PATH) migrate down

.PHONY: dev
dev: ## Run the application in development mode with hot reload
	air

.PHONY: generate
generate: ## Run go generate
	go generate ./...