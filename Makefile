APP_ENV ?= dev
PORT ?= 40001

# Start development server

.PHONY: run
run:
	dotenv -c ${APP_ENV} -- go run cmd/serve/*.go

.PHONY: run_watch
run_watch:
	dotenv -c ${APP_ENV} -- go tool air -build.bin="make run" -build.cmd="/bin/true" -build.include_ext="go,mod"

# Build binaries

.PHONY: build
build: gen
	dotenv -c ${APP_ENV} -- go build -o out/serve ./cmd/serve/*.go

.PHONY: build_for_debug
build_for_debug:
	dotenv -c ${APP_ENV} -- go build -gcflags=all="-N -l" -o out/serve ./cmd/serve/*.go

.PHONY: build_for_debug_watch
build_for_debug_watch:
# this watcher waits 1 second before building, to allow the generators to update Go files (eg: Templ, Gorm).
	dotenv -c ${APP_ENV} -- go tool air -build.bin=out/serve -build.cmd="sleep 1 && make build_for_debug"

.PHONY: build_import_from_moneydance
build_import_from_moneydance: gen
	dotenv -c ${APP_ENV} -- go build -o out/import-from-moneydance ./cmd/import-from-moneydance/*.go

# Generate

.PHONY: gen
gen: gen_templ gen_gorm

# Generate Templ

.PHONY: gen_templ
gen_templ:
	dotenv -c ${APP_ENV} -- go tool templ generate

.PHONY: gen_templ_watch
gen_templ_watch:
	dotenv -c ${APP_ENV} -- go tool templ generate --watch

# Generate Gorm Helpers

.PHONY: gen_gorm
gen_gorm:
	dotenv -c ${APP_ENV} -- go run cmd/dbcodegen/*.go

.PHONY: gen_gorm_watch
gen_gorm_watch:
	dotenv -c ${APP_ENV} -- go tool air -build.bin=/bin/true -build.cmd="make gen_gorm" -build.include_dir="model" -build.exclude_dir="model/query" -build.include_ext="go"

# Generate Styles
.PHONY: gen_css
gen_css:
	npx tailwindcss -i ./view/css/tailwind.css -o ./view/static/css/styles.css

.PHONY: gen_css_watch
gen_css_watch:
	npx tailwindcss -i ./view/css/tailwind.css -o ./view/static/css/styles.css --watch

.PHONY: gen_css_prod
gen_css_prod:
	npx tailwindcss -i ./view/css/tailwind.css -o ./view/static/css/styles.css --minify

# Dev

.PHONY: dev
dev:
	OVERMIND_SKIP_ENV=1 go tool overmind start -f Procfile.dev -p $(PORT)

# Test / Lint / Clean

.PHONY: clean
clean:
	rm -rf out
	find . -name *_templ.go -exec rm {} \;
	find . -name *_templ.txt -exec rm {} \;
	go clean -cache -testcache

.PHONY: test
test:
	dotenv -c test -- go test ./...

.PHONY: ct
ct: clean test

.PHONY: lint
lint: build
	golangci-lint run

.PHONY: check
check: clean gen build lint test

# Database

.PHONY: db_create
db_create:
	dotenv -c ${APP_ENV} -- make _db_create

.PHONY: _db_create
_db_create:
	echo "SELECT 'CREATE DATABASE \"$(PGDATABASE)\"' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$(PGDATABASE)')\gexec" | psql --dbname=postgres

.PHONY: db_migrate
db_migrate:
	dotenv -c ${APP_ENV} -- bin/goose up

.PHONY: db_drop
db_drop:
	dotenv -c ${APP_ENV} -- make _db_drop

.PHONY: _db_drop
_db_drop:
	dotenv -c ${APP_ENV} -- echo "drop database if exists \"$(PGDATABASE)\"" | psql --dbname=postgres

.PHONY: db_seed
db_seed:
	dotenv -c ${APP_ENV} -- go run cmd/dbseed/*.go

.PHONY: psql
psql:
	dotenv -c ${APP_ENV} -- psql

.PHONY: db_full_reset
db_full_reset: db_drop db_create db_migrate db_seed
