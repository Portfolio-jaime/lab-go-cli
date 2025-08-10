# k8s-cli Makefile

.PHONY: build clean install test fmt vet help run-all

# Build the CLI binary
build:
	@echo "Building k8s-cli..."
	go build -o k8s-cli

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f k8s-cli
	go clean

# Install the CLI to GOPATH/bin
install: build
	@echo "Installing k8s-cli..."
	cp k8s-cli $(GOPATH)/bin/

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run all checks
check: fmt vet test

# Run all analysis (requires working cluster)
run-all: build
	@echo "Running complete cluster analysis..."
	./k8s-cli all

# Show version info (requires working cluster)
run-version: build
	@echo "Showing version info..."
	./k8s-cli version

# Show resources (requires working cluster)
run-resources: build
	@echo "Showing resources..."
	./k8s-cli resources

# Show recommendations (requires working cluster)
run-recommend: build
	@echo "Showing recommendations..."
	./k8s-cli recommend

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the k8s-cli binary"
	@echo "  deps          - Install dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  test          - Run tests"
	@echo "  check         - Run fmt, vet, and test"
	@echo "  run-all       - Run complete cluster analysis"
	@echo "  run-version   - Show cluster version"
	@echo "  run-resources - Show cluster resources"
	@echo "  run-recommend - Show recommendations"
	@echo "  help          - Show this help"

# Default target
all: clean build