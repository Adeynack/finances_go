APP_ENV ?= dev
PORT ?= 40001
PGDBPREFIX ?= finances

# Start development server

run:
	dotenv -c ${APP_ENV} -- go run cmd/serve/*.go

run_watch:
	dotenv -c ${APP_ENV} -- air -build.bin="make run" -build.cmd="/bin/true" -build.include_ext="go,mod"

# Build binaries

build: gen
	dotenv -c ${APP_ENV} -- go build -o out/serve ./cmd/serve/*.go

build_for_debug:
	dotenv -c ${APP_ENV} -- go build -gcflags=all="-N -l" -o out/serve ./cmd/serve/*.go

build_for_debug_watch:
# this watcher waits 1 second before building, to allow the generators to update Go files (eg: Templ, Gorm).
	dotenv -c ${APP_ENV} -- air -build.bin=out/serve -build.cmd="sleep 1 && make build_for_debug"

# Generate

gen: gen_templ gen_gorm

# Generate Templ

gen_templ:
	dotenv -c ${APP_ENV} -- templ generate

gen_templ_watch:
	dotenv -c ${APP_ENV} -- templ generate --watch

# Generate Gorm Helpers

gen_gorm:
	dotenv -c ${APP_ENV} -- go run cmd/dbcodegen/*.go

gen_gorm_watch:
	dotenv -c ${APP_ENV} -- air -build.bin=/bin/true -build.cmd="make gen_gorm" -build.include_dir="model" -build.exclude_dir="model/query" -build.include_ext="go"

# Dev

dev:
	OVERMIND_SKIP_ENV=1 overmind start -f Procfile.dev -p $(PORT)

# Test / Lint / Clean

clean:
	rm -rf out
	find . -name *_templ.go -exec rm {} \;
	find . -name *_templ.txt -exec rm {} \;
	go clean -cache -testcache

test:
	dotenv -c test -- go test ./...

ct: clean test

lint: build
	go vet ./...
	bin/staticcheck ./...

check: clean gen build lint test

# Database

db_create:
	dotenv -c ${APP_ENV} -- make _db_create

_db_create:
	echo "SELECT 'CREATE DATABASE \"$(PGDATABASE)\"' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$(PGDATABASE)')\gexec" | psql --dbname=postgres

db_migrate:
	dotenv -c ${APP_ENV} -- bin/goose up

db_drop:
	dotenv -c ${APP_ENV} -- make _db_drop

_db_drop:
	dotenv -c ${APP_ENV} -- echo "drop database if exists \"$(PGDATABASE)\"" | psql --dbname=postgres

db_seed:
	dotenv -c ${APP_ENV} -- go run cmd/dbseed/*.go

psql:
	dotenv -c ${APP_ENV} -- psql

db_full_reset: db_drop db_create db_migrate db_seed

test:
	dotenv -c dev -- env