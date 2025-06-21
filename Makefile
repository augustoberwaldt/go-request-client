.PHONY: test build clean run-examples install

# Build the library
build:
	go build ./...

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Run examples
run-examples:
	go run examples/main.go

# Install dependencies
install:
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	go clean
	rm -f coverage.out

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Generate documentation
docs:
	godoc -http=:6060

# Run benchmarks
bench:
	go test -bench=. ./...

# Check for security vulnerabilities
security:
	gosec ./...

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy

# Create a new release
release:
	@echo "Creating release..."
	@read -p "Enter version (e.g., v1.0.0): " version; \
	git tag $$version; \
	git push origin $$version

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the library"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  run-examples   - Run example code"
	@echo "  install        - Install dependencies"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  docs           - Generate documentation"
	@echo "  bench          - Run benchmarks"
	@echo "  security       - Check for security vulnerabilities"
	@echo "  update-deps    - Update dependencies"
	@echo "  release        - Create a new release"
	@echo "  help           - Show this help" 