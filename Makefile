PORT ?= 40001
DATABASE_NAME ?= finances

# Dev
.PHONY: dev
dev:
	OVERMIND_SKIP_ENV=1 go tool overmind start -f Procfile.dev -p $(PORT)

# Start development server
.PHONY: build_debug
build_debug:
	go tool godotenv go build -gcflags=all="-N -l" -o out/serve ./cmd/serve/*.go

.PHONY: run
run: build_debug
	go tool godotenv out/serve

.PHONY: run_watch
run_watch:
# this watcher waits 1 second before building, to allow the generators to update Go files (eg: Templ, Gorm).
	go tool godotenv go tool air -c air_run_watch.toml

# Build binaries

.PHONY: build
build: gen
	go tool godotenv go build -o bin/serve ./cmd/serve/*.go

.PHONY: gen
gen:
	go tool godotenv go generate ./...

.PHONY: generate_api_code_watch
generate_api_code_watch:
	go tool godotenv go tool air -c air_opai.toml

# Misc

.PHONY: clean
clean:
	go tool godotenv go clean -cache -testcache

.PHONY: test
test:
	go tool godotenv go test ./...

.PHONY: ct
ct: clean test

.PHONY: lint
lint: build
	go vet ./...
	staticcheck ./...
	go tool golangci-lint run

.PHONY: check
check: clean build lint test

# Database

.PHONY: db_create
db_create:
	go tool godotenv make _db_create

.PHONY: _db_create
_db_create:
	echo "SELECT 'CREATE DATABASE \"$(DATABASE_NAME)\"' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$(DATABASE_NAME)')\gexec" | psql --dbname=postgres

.PHONY: db_migrate
db_migrate:
	go tool godotenv go tool goose up
	make db_generate

.PHONY: db_drop
db_drop:
	go tool godotenv make _db_drop

.PHONY: _db_drop
_db_drop:
	echo "drop database if exists \"$(DATABASE_NAME)\"" | psql --dbname=postgres

.PHONY: db_seed
db_seed:
# 	go tool godotenv go run cmd/dbseed/*.go

.PHONY: psql
psql:
	go tool godotenv make _psql

.PHONY: _psql
_psql:
	psql ${DATABASE_URL}

.PHONY: db_full_reset
db_full_reset: db_drop db_create db_migrate db_seed db_generate

.PHONY: db_generate
db_generate:
	go tool godotenv make _db_generate

.PHONY: _db_generate
_db_generate:
	go tool jet -dsn="${DATABASE_URL}" -path=./pkg/repository/gen
