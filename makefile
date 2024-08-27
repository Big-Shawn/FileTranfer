# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Name of the binary
SERVER_BINARY_NAME=filetransfer-server
CLIENT_BINARY_NAME=filetransfer-client

# Platforms
PLATFORMS=darwin linux windows
ARCH=amd64 arm64

all: clean build

build:
	$(GOBUILD) -C server -o $(SERVER_BINARY_NAME) -v
	$(GOBUILD) -C client -o $(CLIENT_BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(SERVER_BINARY_NAME)
	rm -f $(CLIENT_BINARY_NAME)

.PHONY: all build test clean

# Cross compilation
build-all: clean
	@for GOOS in $(PLATFORMS); do \
		for GOARCH in $(ARCH); do \
			$(GOBUILD) -C server -o $(SERVER_BINARY_NAME)-$$GOOS-$$GOARCH -v; \
		done \
	done

	@for GOOS in $(PLATFORMS); do \
    		for GOARCH in $(ARCH); do \
    			$(GOBUILD) -C client -o $(CLIENT_BINARY_NAME)-$$GOOS-$$GOARCH -v; \
    		done \
    	done

.PHONY: build-all