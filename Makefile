# Makefile for wav2multi-lib

.PHONY: help test test-verbose test-coverage build example clean install-deps lint format check

# Default target
help:
	@echo "Available targets:"
	@echo "  make test          - Run tests"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make build         - Build the example"
	@echo "  make example       - Run the example"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make install-deps  - Install dependencies"
	@echo "  make lint          - Run linters"
	@echo "  make format        - Format code"
	@echo "  make check         - Run all checks (format, lint, test)"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build example (without CGO)
build:
	@echo "Building example (without CGO)..."
	cd example && CGO_ENABLED=0 go build -o wav2multi-example main.go

# Build example (with CGO for G.729 support)
build-cgo:
	@echo "Building example (with CGO for G.729)..."
	cd example && CGO_ENABLED=1 CGO_CFLAGS="-I/usr/local/include" CGO_LDFLAGS="-L/usr/local/lib -lbcg729" go build -o wav2multi-example-cgo main.go

# Run example
example:
	@echo "Running example..."
	cd example && go run main.go input.wav

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f example/wav2multi-example
	rm -f example/wav2multi-example-cgo
	rm -f example/output-*.wav
	rm -f example/output_*.g729
	rm -f example/output_*.ulaw
	rm -f example/output_*.alaw
	rm -f example/output_*.slin
	rm -f coverage.out coverage.html

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

# Run linters (requires golangci-lint)
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

# Format code
format:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Run all checks
check: format lint test
	@echo "All checks passed!"

# Install bcg729 on macOS (requires Homebrew or manual installation)
install-bcg729:
	@echo "Installing bcg729..."
	@echo "Cloning bcg729..."
	cd /tmp && rm -rf bcg729 && git clone https://github.com/BelledonneCommunications/bcg729.git
	@echo "Building bcg729..."
	cd /tmp/bcg729 && mkdir -p mybuild && cd mybuild && cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local
	cd /tmp/bcg729/mybuild && make
	@echo "Installing bcg729 (may require sudo)..."
	cd /tmp/bcg729/mybuild && sudo make install || (mkdir -p /tmp/bcg729_install && make install DESTDIR=/tmp/bcg729_install && cp -r /tmp/bcg729_install/usr/local/* /usr/local/)
	@echo "bcg729 installed successfully!"

# Show version info
version:
	@echo "Go version:"
	@go version
	@echo "\nModule info:"
	@go list -m

