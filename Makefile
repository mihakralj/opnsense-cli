# Makefile for opnsense module
# run make help for options

# Variables
VERSION=0.2.0
GO=go
BUILD_DIR=build
BINARY_NAME=opnsense

# Phony Targets
.PHONY: all clean build test install fmt cross-build help

# Default Target
all: clean build

# Display Help
help:
	@echo "Makefile commands for opnsense:"
	@echo "all        Build the binary"
	@echo "clean      Remove build artifacts"
	@echo "build      Build the binary"
	@echo "test       Run tests"
	@echo "install    Install the binary"
	@echo "fmt        Format the source code"
	@echo "cross-build Cross-compile for different platforms"

# Clean up
clean:
	@echo "Cleaning..."
ifeq ($(OS),Windows_NT)
	-@rmdir /s /q $(BUILD_DIR)
else
	-@rm -rf $(BUILD_DIR)
endif
	@echo "Cleaning done"

# Download Dependencies
deps:
	@echo "Downloading dependencies..."
	@$(GO) mod tidy

# Build the binary
build: deps
	@echo "Building..."
ifeq ($(OS),Windows_NT)
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	@$(GO) build -ldflags "-X cmd.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME).exe .
else
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -ldflags "-X cmd.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) .
endif

# Run tests
test: deps
	@echo "Testing..."
	@$(GO) test -v ./...

# Install the binary
install:
	@echo "Installing..."
ifeq ($(OS),Windows_NT)
	@$(GO) install .
else
	@sudo $(GO) install .
endif

# Format the code
fmt:
	@echo "Formatting..."
	@$(GO) fmt ./...

# Cross-compile for different platforms
cross-build: deps
	@echo "Cross-building..."
ifeq ($(OS),Windows_NT)
	@set GOOS=linux&& set GOARCH=amd64&& $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-lnx .
	@echo "Linux binary done"
	@set GOOS=windows&& set GOARCH=amd64&& $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME).exe .
	@echo "Windows binary done"
	@set GOOS=darwin&& set GOARCH=amd64&& $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-mac .
	@echo "MacOS binary done"
	@set GOOS=freebsd&& set GOARCH=amd64&& $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-bsd .
	@echo "FreeBSD binary done"
else
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-lnx .
	@echo "Linux binary done"
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME).exe .
	@echo "Windows binary done"
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-mac .
	@echo "MacOS binary done"
	@GOOS=freebsd GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-bsd .
	@echo "FreeBSD binary done"
endif

