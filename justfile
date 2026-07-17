# Justfile for pathcmds CLI utility

# List all available recipes
default:
	@just --list

# Initialize the module and download dependencies
init:
	go mod init pathcmds || true
	go mod tidy
	go mod vendor

# Compile the pathcmds executable into bin/pathcmds
build:
	mkdir -p bin
	go build -mod=vendor -o bin/pathcmds main.go

# Install the executable to the GOPATH bin directory
install:
	go install -mod=vendor

# Uninstall the executable from Go bin directories
uninstall:
	go clean -i || true
	rm -f "$(go env GOBIN)/pathcmds" || true
	rm -f "$(go env GOPATH)/bin/pathcmds" || true
	rm -f "$HOME/go/bin/pathcmds" || true

# Clean up build artifacts and Go build caches
clean:
	rm -rf bin
	go clean -cache

# Run the utility immediately with optional CLI flags (e.g. just run "--system --page")
run flags="":
	go run -mod=vendor main.go {{flags}}

# Run all CI checks: format check, vet (lint), typecheck, and unit tests
ci: fmt-check lint typecheck test

# Verify all Go source files are formatted correctly
fmt-check:
	@if [ -n "$(gofmt -l cmd pkg main.go)" ]; then \
		echo "The following files are not formatted correctly. Please run 'go fmt ./...':"; \
		gofmt -l cmd pkg main.go; \
		exit 1; \
	fi

# Lint the codebase using go vet
lint:
	go vet -mod=vendor ./...

# Verify that the project compiles correctly without writing the binary
typecheck:
	go build -mod=vendor -o /dev/null main.go

# Run the unit tests with coverage
test:
	go test -mod=vendor -v -cover ./...
