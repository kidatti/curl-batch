# Variables
BINARY_NAME=curl-batch
VERSION?=latest
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION}"

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# Test
.PHONY: test
test:
	go test ./...

# Clean
.PHONY: clean
clean:
	rm -f ${BINARY_NAME}
	rm -rf dist/

# Build for all platforms
.PHONY: build-all
build-all: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-arm64.exe .

# Build for Linux
.PHONY: build-linux
build-linux:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 .

# Build for macOS
.PHONY: build-darwin
build-darwin:
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .

# Build for Windows
.PHONY: build-windows
build-windows:
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-arm64.exe .

# Run with sample data
.PHONY: run-sample
run-sample: build
	./${BINARY_NAME} -curl sample/curl.txt -csv sample/hogehoge.csv -output results.txt

# Install (copy to /usr/local/bin)
.PHONY: install
install: build
	sudo cp ${BINARY_NAME} /usr/local/bin/

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	elif [ -f "$$(go env GOPATH)/bin/golangci-lint" ]; then \
		$$(go env GOPATH)/bin/golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		$$(go env GOPATH)/bin/golangci-lint run; \
	fi

# Create release packages
.PHONY: package
package: build-all
	mkdir -p dist/packages
	cd dist && tar -czf packages/${BINARY_NAME}-linux-amd64.tar.gz ${BINARY_NAME}-linux-amd64
	cd dist && tar -czf packages/${BINARY_NAME}-linux-arm64.tar.gz ${BINARY_NAME}-linux-arm64
	cd dist && tar -czf packages/${BINARY_NAME}-darwin-amd64.tar.gz ${BINARY_NAME}-darwin-amd64
	cd dist && tar -czf packages/${BINARY_NAME}-darwin-arm64.tar.gz ${BINARY_NAME}-darwin-arm64
	cd dist && zip packages/${BINARY_NAME}-windows-amd64.zip ${BINARY_NAME}-windows-amd64.exe
	cd dist && zip packages/${BINARY_NAME}-windows-arm64.zip ${BINARY_NAME}-windows-arm64.exe

# Create and push git tag for release
.PHONY: tag
tag:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "latest" ]; then \
		echo "Please specify VERSION: make tag VERSION=v1.0.0"; \
		exit 1; \
	fi
	git tag $(VERSION)
	git push origin $(VERSION)
	@echo "Tagged and pushed $(VERSION). GitHub Actions will create the release."

# Remove git tag (local and remote)
.PHONY: untag
untag:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "latest" ]; then \
		echo "Please specify VERSION: make untag VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Removing tag $(VERSION)..."
	git tag -d $(VERSION) || echo "Local tag $(VERSION) not found"
	git push origin :refs/tags/$(VERSION) || echo "Remote tag $(VERSION) not found"
	@echo "Tag $(VERSION) removed from local and remote repositories"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms"
	@echo "  build-linux  - Build for Linux (amd64, arm64)"
	@echo "  build-darwin - Build for macOS (amd64, arm64)"
	@echo "  build-windows- Build for Windows (amd64, arm64)"
	@echo "  package      - Create release packages (tar.gz, zip)"
	@echo "  tag          - Create and push git tag (use VERSION=v1.0.0)"
	@echo "  untag        - Remove git tag locally and remotely (use VERSION=v1.0.0)"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  run-sample   - Build and run with sample data"
	@echo "  install      - Install binary to /usr/local/bin"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Release workflow:"
	@echo "  1. make tag VERSION=v1.0.0"
	@echo "  2. GitHub Actions will build and release automatically"
	@echo ""
	@echo "Usage example:"
	@echo "  ./curl-batch -curl sample/curl.txt -csv sample/hogehoge.csv -output results.txt -sleep 1000"