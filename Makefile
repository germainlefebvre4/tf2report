# tf2report Makefile
# Builds the Terraform plan analysis CLI tool

# Variables
BINARY_NAME := tf2report
CMD_DIR := cmd/tf2report
BUILD_DIR := bin
GO := go
GOFLAGS := -v
LDFLAGS := -s -w
PREFIX ?= /usr/local
DESTDIR ?=

# Go commands
GOCMD := $(GO)
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOFMT := $(GOCMD) fmt
GOMOD := $(GOCMD) mod

# Build variables
VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
BUILD_LDFLAGS := $(LDFLAGS) -X 'main.appVersion=$(VERSION)' -X 'main.buildCommit=$(COMMIT)' -X 'main.buildTime=$(BUILD_TIME)'

# Default target
.PHONY: build install test clean run help

# Build the application
.PHONY: build
build:
	@echo "Building $(BUILD_DIR)/$(BINARY_NAME)..."
	$(GO) build -ldflags="$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) cmd/tf2report/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install the binary to system path
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
# 	sudo install -d $(DESTDIR)$(PREFIX)/bin
	sudo install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(DESTDIR)$(PREFIX)/bin/
	@echo "Installation complete. Run 'tf2report --help' to get started."

# Uninstall binary from system
.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(BINARY_NAME)

# Format code
.PHONY: fmt
fmt:
	$(GOFMT) ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Verify dependencies
.PHONY: verify
verify:
	$(GOMOD) verify

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download

# Run the application
.PHONY: run
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Build for multiple platforms
.PHONY: release
release: clean
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(GOFLAGS) -ldflags "$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(GOFLAGS) -ldflags "$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(GOFLAGS) -ldflags "$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(GOFLAGS) -ldflags "$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(GOFLAGS) -ldflags "$(BUILD_LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)

# Help target
.PHONY: help
help:
	@echo "tf2report - Terraform Plan Analysis Tool"
	@echo ""
	@echo "Available targets:"
	@echo "  all          - Build the binary (default)"
	@echo "  build        - Build the binary"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean        - Remove build artifacts"
	@echo "  install      - Install binary to $(PREFIX)/bin"
	@echo "  uninstall    - Remove binary from $(PREFIX)/bin"
	@echo "  fmt          - Format Go code"
	@echo "  tidy         - Tidy Go modules"
	@echo "  verify       - Verify dependencies"
	@echo "  deps         - Download dependencies"
	@echo "  run          - Build and run the application"
	@echo "  release      - Build for multiple platforms"
	@echo "  help         - Show this help message"
