.PHONY: all build test lint ci clean

# Binary name
BINARY_NAME=main

# Build the application
build:
	@echo "Building..."
	go build -o $(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Run all CI checks (lint, test, build)
ci: lint test build
	@echo "CI checks completed successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)

# Run the application
run: build
	@echo "Running application..."
	./$(BINARY_NAME)
