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

.PHONY: e2e
e2e:
	@echo "Ejecutando pruebas end-to-end..."
	./test/e2e/test_k8s_cli.sh
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

# NEW: Run real-time metrics (requires working cluster)
run-metrics: build
	@echo "Showing real-time metrics..."
	./k8s-cli metrics --nodes --pods --utilization

# NEW: Run cost analysis (requires working cluster)
run-cost: build
	@echo "Showing cost analysis..."
	./k8s-cli cost

# NEW: Run workload analysis (requires working cluster)
run-workload: build
	@echo "Showing workload health analysis..."
	./k8s-cli workload

# NEW: Run logs analysis (requires working cluster)
run-logs: build
	@echo "Showing logs and events analysis..."
	./k8s-cli logs --critical --patterns

# NEW: Export data (requires working cluster)
run-export: build
	@echo "Exporting cluster data..."
	./k8s-cli export --format json --costs --metrics

# NEW: Demo script
demo: build
	@echo "Running demo of new features..."
	./examples/demo_new_features.sh

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
	@echo "  run-metrics   - Show real-time metrics (NEW)"
	@echo "  run-cost      - Show cost analysis (NEW)"
	@echo "  run-workload  - Show workload health (NEW)"
	@echo "  run-logs      - Show logs analysis (NEW)"
	@echo "  run-export    - Export cluster data (NEW)"
	@echo "  demo          - Run demo of new features (NEW)"
	@echo "  help          - Show this help"

# Default target
all: clean build