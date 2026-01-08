.PHONY: build run clean fmt test

# Binary name
BINARY=herzog-drei

# Build flags for macOS
# raylib-go uses CGO and requires specific frameworks on macOS
CGO_ENABLED=1
CGO_LDFLAGS=-framework CoreVideo -framework IOKit -framework Cocoa -framework GLUT -framework OpenGL

build:
	CGO_ENABLED=$(CGO_ENABLED) go build -o $(BINARY) .

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	go clean

fmt:
	go fmt ./...

test:
	go test ./...

# Development mode with race detector
dev:
	CGO_ENABLED=$(CGO_ENABLED) go build -race -o $(BINARY) .
	./$(BINARY)

# Build for release (optimized)
release:
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags="-s -w" -o $(BINARY) .
