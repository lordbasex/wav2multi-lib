# Makefile for wav2multi-lib

.PHONY: help test test-verbose test-coverage build example clean install-deps lint format check tag tag-push release deploy tag-delete tag-list

# Default target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Development:"
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
	@echo ""
	@echo "Release & Deployment:"
	@echo "  make tag           - Create a new version tag (interactive)"
	@echo "  make tag-push      - Push tags to remote"
	@echo "  make release       - Create and push a new release tag"
	@echo "  make deploy        - Full deployment (check + tag + push)"
	@echo "  make tag-list      - List all tags"
	@echo "  make tag-delete    - Delete a tag (local and remote)"
	@echo ""
	@echo "Setup:"
	@echo "  make install-bcg729 - Install bcg729 library (G.729 codec)"
	@echo "  make version       - Show Go and module version"

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

# ============================================================================
# Release & Deployment
# ============================================================================

# Create a new version tag (interactive)
tag:
	@echo "Creating a new version tag..."
	@echo ""
	@echo "Current tags:"
	@git tag -l | tail -5 || echo "  (no tags yet)"
	@echo ""
	@read -p "Enter new version (e.g., v1.0.0): " version; \
	if [ -z "$$version" ]; then \
		echo "❌ Version cannot be empty"; \
		exit 1; \
	fi; \
	if ! echo "$$version" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "❌ Version must follow format: vX.Y.Z (e.g., v1.0.0)"; \
		exit 1; \
	fi; \
	if git tag -l | grep -q "^$$version$$"; then \
		echo "❌ Tag $$version already exists"; \
		exit 1; \
	fi; \
	read -p "Enter release message: " message; \
	if [ -z "$$message" ]; then \
		message="Release $$version"; \
	fi; \
	git tag -a "$$version" -m "$$message"; \
	echo "✅ Tag $$version created successfully!"; \
	echo ""; \
	echo "Next steps:"; \
	echo "  1. Push tag:    make tag-push"; \
	echo "  2. Or delete:   make tag-delete VERSION=$$version"

# Push tags to remote
tag-push:
	@echo "Pushing tags to remote..."
	@if [ -z "$$(git remote)" ]; then \
		echo "❌ No remote repository configured"; \
		exit 1; \
	fi; \
	git push origin --tags
	@echo "✅ Tags pushed successfully!"
	@echo ""
	@echo "🎉 Release published! Users can now install with:"
	@echo "   go get github.com/lordbasex/wav2multi-lib@$$(git describe --tags --abbrev=0)"

# Create and push a new release tag (non-interactive)
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required"; \
		echo "Usage: make release VERSION=v1.0.0 MESSAGE='Release message'"; \
		exit 1; \
	fi; \
	if ! echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "❌ VERSION must follow format: vX.Y.Z (e.g., v1.0.0)"; \
		exit 1; \
	fi; \
	if git tag -l | grep -q "^$(VERSION)$$"; then \
		echo "❌ Tag $(VERSION) already exists"; \
		exit 1; \
	fi; \
	MSG="$${MESSAGE:-Release $(VERSION)}"; \
	git tag -a "$(VERSION)" -m "$$MSG"; \
	echo "✅ Tag $(VERSION) created"; \
	git push origin --tags; \
	echo "✅ Tag pushed to remote"; \
	echo ""; \
	echo "🎉 Release $(VERSION) published!"; \
	echo "   go get github.com/lordbasex/wav2multi-lib@$(VERSION)"

# Full deployment: checks + tag + push
deploy: check
	@echo ""
	@echo "╔═══════════════════════════════════════════════════════════╗"
	@echo "║          🚀 DEPLOYMENT PROCESS                            ║"
	@echo "╚═══════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "✅ All checks passed!"
	@echo ""
	@echo "Current tags:"
	@git tag -l | tail -5 || echo "  (no tags yet)"
	@echo ""
	@read -p "Enter new version (e.g., v1.0.0): " version; \
	if [ -z "$$version" ]; then \
		echo "❌ Version cannot be empty"; \
		exit 1; \
	fi; \
	if ! echo "$$version" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "❌ Version must follow format: vX.Y.Z (e.g., v1.0.0)"; \
		exit 1; \
	fi; \
	if git tag -l | grep -q "^$$version$$"; then \
		echo "❌ Tag $$version already exists"; \
		exit 1; \
	fi; \
	read -p "Enter release message (optional): " message; \
	if [ -z "$$message" ]; then \
		message="Release $$version"; \
	fi; \
	echo ""; \
	echo "📋 Deployment Summary:"; \
	echo "   Version: $$version"; \
	echo "   Message: $$message"; \
	echo "   Remote:  $$(git remote get-url origin 2>/dev/null || echo 'not configured')"; \
	echo ""; \
	read -p "Proceed with deployment? [y/N] " confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "❌ Deployment cancelled"; \
		exit 1; \
	fi; \
	echo ""; \
	echo "Creating tag..."; \
	git tag -a "$$version" -m "$$message"; \
	echo "✅ Tag $$version created"; \
	echo ""; \
	echo "Pushing to remote..."; \
	git push origin --tags; \
	echo "✅ Tag pushed to remote"; \
	echo ""; \
	echo "╔═══════════════════════════════════════════════════════════╗"; \
	echo "║          🎉 DEPLOYMENT SUCCESSFUL!                        ║"; \
	echo "╚═══════════════════════════════════════════════════════════╝"; \
	echo ""; \
	echo "📦 Users can now install with:"; \
	echo "   go get github.com/lordbasex/wav2multi-lib@$$version"; \
	echo ""; \
	echo "🔗 Next steps:"; \
	echo "   1. Create GitHub Release (optional):"; \
	echo "      https://github.com/lordbasex/wav2multi-lib/releases/new?tag=$$version"; \
	echo "   2. Update CHANGELOG.md if needed"; \
	echo "   3. Announce the release"

# List all tags
tag-list:
	@echo "📋 All tags:"
	@git tag -l -n1 || echo "  (no tags yet)"

# Delete a tag (local and remote)
tag-delete:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required"; \
		echo "Usage: make tag-delete VERSION=v1.0.0"; \
		exit 1; \
	fi; \
	echo "Deleting tag $(VERSION)..."; \
	git tag -d "$(VERSION)" 2>/dev/null || echo "  Local tag not found"; \
	git push origin :refs/tags/$(VERSION) 2>/dev/null || echo "  Remote tag not found"; \
	echo "✅ Tag $(VERSION) deleted (if it existed)"

