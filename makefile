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
	rm -f ./server/$(SERVER_BINARY_NAME)*
	rm -f ./client/$(CLIENT_BINARY_NAME)*

.PHONY: all build test clean

# Cross compilation
build-all: clean
	@for GOOS in $(PLATFORMS); do \
		for GOARCH in $(ARCH); do \
		    echo "Server Building for $$GOOS/$$GOARCH..."; \
		    if [[ "$$GOOS" == "windows" ]]; then \
			GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) -C server -o $(SERVER_BINARY_NAME)-$$GOOS-$$GOARCH.exe -v; \
		    else \
			GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) -C server -o $(SERVER_BINARY_NAME)-$$GOOS-$$GOARCH -v; \
            fi \
		done \
	done

	@for GOOS in $(PLATFORMS); do \
    		for GOARCH in $(ARCH); do \
    		    echo "Client Building for $$GOOS/$$GOARCH..."; \
    		    if [[ "$$GOOS" == "windows" ]]; then \
    			GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) -C client -o $(CLIENT_BINARY_NAME)-$$GOOS-$$GOARCH.exe -v; \
    			else \
    			GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) -C client -o $(CLIENT_BINARY_NAME)-$$GOOS-$$GOARCH -v; \
    		    fi \
    		done \
    	done

.PHONY: build-all