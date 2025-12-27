.PHONY: run build dev install migrate clean help

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o bin/api.exe cmd/api/main.go

# Run with hot reload (requires air)
dev:
	air

# Install dependencies
install:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf tmp/
	rm -rf bin/

# Install air for hot reload
install-air:
	go install github.com/air-verse/air@latest

# Show help
help:
	@echo "Available commands:"
	@echo "  make run         - Run the application"
	@echo "  make build       - Build the application"
	@echo "  make dev         - Run with hot reload (requires air)"
	@echo "  make install     - Install dependencies"
	@echo "  make install-air - Install air for hot reload"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make help        - Show this help message"
