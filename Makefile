PORT ?= 40001

# Dev
.PHONY: dev
dev: tools
	OVERMIND_SKIP_ENV=1 overmind start -f Procfile.dev -p $(PORT)

.PHONY: tools
tools:
	@which -s air || go install github.com/cosmtrek/air@latest
	@which -s staticcheck || go install honnef.co/go/tools/cmd/staticcheck@latest
	@which -s golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2
	@which -s godotenv || go install github.com/joho/godotenv/cmd/godotenv@latest

# Start development server
.PHONY: build_debug
build_debug:
	godotenv go build -gcflags=all="-N -l" -o out/serve ./cmd/serve/*.go

.PHONY: run
run: build_debug
	godotenv out/serve

.PHONY: run_watch
run_watch:
# this watcher waits 1 second before building, to allow the generators to update Go files (eg: Templ, Gorm).
	godotenv air -c air_run_watch.toml

# Build binaries

.PHONY: build
build: gen
	godotenv go build -o bin/serve ./cmd/serve/*.go

.PHONY: gen
gen:
	godotenv go generate ./...

.PHONY: generate_api_code_watch
generate_api_code_watch:
	godotenv air -c air_opai.toml

# Misc

.PHONY: clean
clean:
	godotenv go clean -cache -testcache

.PHONY: test
test:
	godotenv go test ./...

.PHONY: ct
ct: clean test

.PHONY: lint
lint: tools build
	go vet ./...
	staticcheck ./...
	golangci-lint run

.PHONY: check
check: clean build lint test
