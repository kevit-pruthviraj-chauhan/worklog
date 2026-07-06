.PHONY: build build-all clean install release

# Build for current platform
build:
	@echo "Building worklog for $(shell go env GOOS)/$(shell go env GOARCH)..."
	go build -o worklog ./cmd

# Build for multiple platforms
build-all:
	@echo "Building worklog for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o worklog-linux-amd64 ./cmd
	GOOS=linux GOARCH=arm64 go build -o worklog-linux-arm64 ./cmd
	GOOS=darwin GOARCH=amd64 go build -o worklog-darwin-amd64 ./cmd
	GOOS=darwin GOARCH=arm64 go build -o worklog-darwin-arm64 ./cmd
	@echo "Built binaries:"
	@ls -lh worklog-*

# Install to /usr/local/bin
install: build
	@echo "Installing worklog to /usr/local/bin..."
	sudo mv worklog /usr/local/bin/
	@echo "Installation complete!"

# Install from current build
install-local:
	@echo "Installing worklog to /usr/local/bin..."
	sudo mv worklog /usr/local/bin/
	@echo "Installation complete!"

# Create a new release tag
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release $(VERSION)..."
	git tag $(VERSION)
	git push origin $(VERSION)
	@echo "Release created! Check GitHub Actions for build progress."

# Clean build artifacts
clean:
	rm -f worklog worklog-*
	@echo "Cleaned build artifacts"
