APP_NAME := bookmarks-server
DIST_DIR := dist

.PHONY: all build clean release-arm64

all: build

# Standard build for current OS
build:
	go build -o $(APP_NAME) cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf $(DIST_DIR) $(APP_NAME)

# Cross-compilation for Linux ARM64 and archive creation
release-arm64: clean
	mkdir -p $(DIST_DIR)/$(APP_NAME)
	# Compile for Linux ARM64 (CGO_ENABLED=0 for pure Go sqlite)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(DIST_DIR)/$(APP_NAME)/$(APP_NAME) cmd/server/main.go
	# Copy templates
	cp -r web $(DIST_DIR)/$(APP_NAME)/
	# Create archive
	cd $(DIST_DIR) && tar -czvf $(APP_NAME)-linux-arm64.tar.gz $(APP_NAME)
	@echo "Archive ready: $(DIST_DIR)/$(APP_NAME)-linux-arm64.tar.gz"
