PORT ?= 40001

# Start development server
run:
	go run cmd/serve/*.go

run_watch:
	air -build.bin="make run" -build.cmd="/bin/true" -build.include_ext="go,mod"

# Build binaries
build:
	go build -o out/serve ./cmd/serve/*.go

build_for_debug:
	go build -gcflags=all="-N -l" -o out/serve ./cmd/serve/*.go

build_for_debug_watch:
	air -build.bin=out/serve -build.cmd "make build_for_debug"

# Generate

gen_templ:
	templ generate

gen_templ_watch:
	templ generate --watch

gen: gen_templ

# Dev
dev:
	overmind start -f Procfile.dev -p $(PORT)

# Misc

clean:
	rm -rf out
	find . -name *_templ.go -exec rm {} \;
	find . -name *_templ.txt -exec rm {} \;
	go clean -cache -testcache

test:
	go test ./...

ct: clean test

lint: build
	go vet ./...
	staticcheck ./...

check: clean gen build lint test
