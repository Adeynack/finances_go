PORT ?= 40001
APP_ENV ?= dev
PGDBPREFIX ?= finances

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
# this watcher waits 1 second before building, to allow the generators to update Go files (eg: Templ, Gorm).
	air -build.bin=out/serve -build.cmd="sleep 1 && make build_for_debug"

# Generate

gen: gen_templ gen_gorm

# Generate Templ

gen_templ:
	templ generate

gen_templ_watch:
	templ generate --watch

# Generate Gorm Helpers

gen_gorm:
	go run cmd/ops/*.go db:codegen

gen_gorm_watch:
	air -build.bin=/bin/true -build.cmd="make gen_gorm" -build.include_dir="model" -build.exclude_dir="model/query" -build.include_ext="go"

# Dev

dev:
	overmind start -f Procfile.dev -p $(PORT)

# Test / Lint / Clean

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

# Database

db_create:
	echo "SELECT 'CREATE DATABASE \"$(PGDBPREFIX)_dev\"' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'finances')\gexec" | psql
	echo "SELECT 'CREATE DATABASE \"$(PGDBPREFIX)_test\"' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'finances')\gexec" | psql

db_migrate:
	PGDB="$(PGDBPREFIX)_$(APP_ENV)" bin/goose up

db_drop:
	echo "drop database if exists \"$(PGDBPREFIX)_dev\"" | psql
	echo "drop database if exists \"$(PGDBPREFIX)_test\"" | psql

db_seed:
	echo "TODO: Seed"

db_full_reset: db_drop db_create db_migrate db_seed
