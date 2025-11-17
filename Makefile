.PHONY: build test install clean fmt lint run help

# Build binary
build:
	@echo "Building copilot-research..."
	@go build -o copilot-research -ldflags="-s -w"
	@echo "✅ Build complete"

# Run tests with coverage
test:
	@echo "Running tests..."
	@go test ./... -v -cover -coverprofile=coverage.txt
	@echo "✅ Tests complete"

# Install to GOPATH
install:
	@echo "Installing..."
	@go install
	@echo "✅ Installed to $(shell go env GOPATH)/bin/copilot-research"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f copilot-research coverage.txt
	@rm -rf tmp/*
	@echo "✅ Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .
	@echo "✅ Format complete"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "✅ Lint complete"

# Run directly
run:
	@go run main.go $(ARGS)

# Show help
help:
	@echo "Copilot Research - Makefile Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  build    - Build binary"
	@echo "  test     - Run tests with coverage"
	@echo "  install  - Install to GOPATH"
	@echo "  clean    - Remove build artifacts"
	@echo "  fmt      - Format code"
	@echo "  lint     - Run linter (requires golangci-lint)"
	@echo "  run      - Run directly (use ARGS='query' to pass arguments)"
	@echo "  help     - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make test"
	@echo "  make run ARGS='\"test query\"'"
