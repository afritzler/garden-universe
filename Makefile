GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=garden-universe
BINARY_LINUX=$(BINARY_NAME)_linux_amd64
BINARY_DARWIN=$(BINARY_NAME)_darwin_amd64
BINARY_WIN=$(BINARY_NAME)_win_amd64.exe
IMAGE=afritzler/garden-universe
TAG=latest

all: test build

.PHONY: build
build: build-web
		@CGO_ENABLED=0 GO111MODULE=on $(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: test
test:
		$(GOTEST) -v ./...

.PHONY: clean
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_LINUX)
		rm -f $(BINARY_DARWIN)
		rm -f $(BINARY_WIN)


.PHONY: run
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

.PHONY: dep
deps:
		$(GOGET) -u github.com/rakyll/statik

.PHONY: build-web
build-web:
		statik -src=$(PWD)/web/

.PHONY: revendor
revendor:
		@GO111MODULE=on go mod vendor
		@GO111MODULE=on go mod tidy

# Cross compilation
.PHONY: build-linux
cross: build-linux build-darwin build-win

.PHONY: build-linux
build-linux: build-web
		@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on $(GOBUILD) -o $(BINARY_LINUX) -v

.PHONY: build-darwin
build-darwin: build-web
		@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on $(GOBUILD) -o $(BINARY_DARWIN) -v

.PHONY: build-win
build-win: build-web
		@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on $(GOBUILD) -o $(BINARY_WIN) -v

.PHONY: docker-build
docker-build:
		docker build -t $(IMAGE):$(TAG) .
