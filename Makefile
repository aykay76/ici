# Build the ici binary
build:
	go build -o ici ./cmd/ici

# Install to /usr/local/bin
install: build
	sudo mv ici /usr/local/bin/

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Download dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -f ici
	rm -f coverage.out

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run example workflow
example: build
	./ici parse .github/workflows/example.yml

.PHONY: build install test test-coverage deps clean lint fmt example
