# Start development server
.PHONY: run
run: gen
	go run cmd/serve/*.go

# Build binaries
.PHONY: build
build: gen
	go build -o bin/serve ./cmd/serve/*.go

.PHONY: gen
gen:
	go generate ./...

# Misc

.PHONY: clean
clean:
	go clean -cache -testcache

.PHONY: test
test:
	go test ./...

.PHONY: ct
ct: clean test

.PHONY: lint
lint: build
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

.PHONY: check
check: clean build lint test
