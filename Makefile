export

GO_PACKAGES := $(shell go list ./... | grep -v /vendor/)

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: run
run: build ## Run the binary
	.bin/vertex

.PHONY: run-server
run-server: ## Run the server
	go run ./cmd/server/main.go

.PHONY: build
build: ## Build the binary
	go build -ldflags="-extldflags=-static" -o .bin/vertex ./cli/vertex.go

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

.PHONY: optimize
optimize: ## Optimize structs 
	betteralign -apply ./...

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -f docker/Dockerfile -t vertex:latest .

docker-run: ### Build & Run Vertex with docker
	docker run --name vertex -p 6380:6380 --rm  --detach vertex:latest


