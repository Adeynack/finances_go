PORT ?= 40001

.PHONY: ensure_tools
ensure_tools:
	which air || go install github.com/cosmtrek/air@latest
	which staticcheck || go install honnef.co/go/tools/cmd/staticcheck@latest

# Start development server
.PHONY: run
run: gen
	go run cmd/serve/*.go

.PHONY: build_for_debug
build_for_debug:
	go build -gcflags=all="-N -l" -o out/serve ./cmd/serve/*.go

.PHONY: build_for_debug_watch
build_for_debug_watch:
# this watcher waits 1 second before building, to allow the generators to update Go files (eg: Templ, Gorm).
	air -build.bin=out/serve -build.cmd="sleep 1 && make build_for_debug" -build.include_ext="go" -build.exclude_dir="node_modules,out,tmp"

# Build binaries

.PHONY: build
build: gen
	go build -o bin/serve ./cmd/serve/*.go

.PHONY: gen
gen:
	go generate ./...

.PHONY: generate_api_code_watch
generate_api_code_watch:
	air -build.bin=true -build.cmd="go generate pkg/api/package.go" -build.include_dir="pkg/api" -build.include_ext="yaml"

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
lint: ensure_tools build
	go vet ./...
	staticcheck ./...

.PHONY: check
check: clean build lint test

# Dev
.PHONY: dev
dev:
	OVERMIND_SKIP_ENV=1 overmind start -f Procfile.dev -p $(PORT)
