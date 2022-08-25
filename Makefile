BINARY_NAME=garden-universe

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Image URL to use all building/pushing image targets
IMG ?= garden-universe:latest

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

all: test build ## Test and build.

.PHONY: build
build: build-web ## Build binary.
	go build -o $(BINARY_NAME) -v

lint: ## Run golangci-lint against code.
	golangci-lint run ./...

check: lint test

.PHONY: test
test: deps ## Run tests.
	go test -v ./...

.PHONY: clean
clean: ## Remove build artefacts.
	go clean
	rm -f $(BINARY_NAME)

.PHONY: run
run: ## Run locally.
	go run -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

.PHONY: deps
deps: ## Get dependencies.
	go get -u github.com/rakyll/statik

.PHONY: build-web
build-web: deps ## Regenerate web content.
	statik -f -src=$(PWD)/web/

.PHONY: docker-build
docker-build: ## Build docker image.
	docker build -t ${IMG} .
