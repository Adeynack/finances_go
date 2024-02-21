# Start development server
run:
	go run cmd/serve/*.go

# Build binaries
build:
	go build -o bin/serve ./cmd/serve/*.go

# Misc

clean:
	go clean -cache -testcache

test:
	go test ./...

ct: clean test

lint: build
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

check: clean build lint test
