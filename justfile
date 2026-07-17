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
