export

GO_PACKAGES := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
GO_FILES_NOVENDOR := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v /bin/)

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Setup the project
	go mod download

.PHONY: build
build: build-cli build-server

.PHONY: build-cli
build-cli: ## Build the CLI binary
	go build -ldflags="-extldflags=-static" -o .bin/vertex ./internal/cli/cli.go

.PHONY: build-server
build-server: ## Build the server binary
	go build -ldflags="-extldflags=-static" -o .bin/server ./server/main.go

.PHONY: run-cli
run-cli: build-cli ## Run the CLI binary
	.bin/vertex

.PHONY: run-cli-dev
run-cli-dev: ## Run the CLI for development
	go run ./internal/cli/cli.go

.PHONY: run-server
run-server: build-server ## Run the server binary
	.bin/server

.PHONY: run-server-dev
run-server-dev: ## Run the server for development
	go run ./server/main.go

.PHONY: test
test: ## Run the unit test, make test ARGS=location
	@if [ "$(ARGS)" = "" ]; then \
		echo "Running all tests"; \
		go test ./...; \
	else \
		echo "Running tests in $(ARGS)"; \
		go test $(ARGS); \
	fi

.PHONY: lint
lint: ## Lint Go files
	golangci-lint run ./...

.PHONY: optimize-structs
optimize: ## Optimize structs 
	betteralign -apply ./...

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -f docker/Dockerfile -t vertex:latest .

.PHONY: docker-run
docker-run: ## Build & Run Vertex with docker
	docker run --name vertex -p 6380:6380 --rm --detach vertex:latest

