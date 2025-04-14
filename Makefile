APP_NAME := powerdns-webui
SRC_DIR := .
BUILD_DIR := build
BIN := $(BUILD_DIR)/$(APP_NAME)

GO := go
GOTEST := $(GO) test
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOMOD := $(GO) mod

.PHONY: all build run test clean tidy

all: build

build:
    @echo "Building the application..."
    @mkdir -p $(BUILD_DIR)
    $(GOBUILD) -o $(BIN)

run: build
    @echo "Running the application..."
    $(BIN)

test:
    @echo "Running tests..."
    $(GOTEST) ./...

clean:
    @echo "Cleaning up..."
    $(GOCLEAN)
    rm -rf $(BUILD_DIR)

tidy:
    @echo "Tidying up dependencies..."
    $(GOMOD) tidy
