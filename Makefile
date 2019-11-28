EXECUTABLE=ytpf
RELEASE_DIR=./bin
WINDOWS=$(RELEASE_DIR)/$(EXECUTABLE)_windows_amd64.exe
LINUX=$(RELEASE_DIR)/$(EXECUTABLE)_linux_amd64
DARWIN=$(RELEASE_DIR)/$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always --long --dirty)

all: test build ## Build and run tests

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -i -v -o $(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)" ./cmd/ytpf.go

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -i -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)" ./cmd/ytpf.go

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -i -v -o $(DARWIN) -ldflags="-s -w -X main.version=$(VERSION)" ./cmd/ytpf.go

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

build: windows linux darwin ## Build binaries
	@echo version: $(VERSION)


test: ## Run unit tests
	go test ./...

clean: ## Remove previous build
	rm -rf $(RELEASE_DIR)

.PHONY: all test clean
