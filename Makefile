MAIN_FILE := main.go
BINARY_NAME := APIProbe
EXT := $(if $(filter windows,$(GOOS))$(filter Windows_NT,$(OS)),.exe)
ZIP_TOOL := "/c/Program Files/7-Zip/7z.exe"
BUILD_DIR := dist
COVERAGE_DIR := coverage
PLATFORMS := windows/amd64 windows/arm64 darwin/amd64 darwin/arm64

.PHONY: help run build zip cross release lint test cover tidy clean format

help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Run the application
	go run $(MAIN_FILE)

build: ## Build application for current platform
	go build -o $(BINARY_NAME)$(EXT) $(MAIN_FILE)

zip: build ## Create zip archive of the application
	@rm -f $(BINARY_NAME).zip
	$(ZIP_TOOL) a $(BINARY_NAME).zip $(BINARY_NAME)$(EXT)

cross: ## Cross-compile for different platforms
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%%/*}; \
		GOARCH=$${platform##*/}; \
		EXT=$$( [ "$$GOOS" = "windows" ] && echo ".exe" || echo "" ); \
		OUTPUT=$(BUILD_DIR)/$(BINARY_NAME)-$${GOOS}-$${GOARCH}$$EXT; \
		echo "==> Building $$OUTPUT"; \
		GOOS=$$GOOS GOARCH=$$GOARCH CGO_ENABLED=0 go build -o $$OUTPUT $(MAIN_FILE) || echo "❌ Build failed for $$platform"; \
	done

release: cross ## Zip all compiled platforms binaries
	@cd $(BUILD_DIR) && $(ZIP_TOOL) a $(BINARY_NAME)-release.zip *

lint: ## Run linter
	golangci-lint run ./...

test: ## Run tests
	go test -v ./...

cover: ## Generate coverage report
	@mkdir -p $(COVERAGE_DIR)
	go test -v ./... -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o ./$(COVERAGE_DIR)/coverage.html

tidy: ## Ensure go.mod dependencies are tidy
	go mod tidy

clean: ## Remove build artifacts
	@rm -f $(BINARY_NAME)$(EXT) $(BINARY_NAME).zip
	@rm -rf ./$(BUILD_DIR)
	@rm -rf ./$(COVERAGE_DIR)

format: ## Format the file (LF instead of CRLF) using gofmt
	go fmt ./...
