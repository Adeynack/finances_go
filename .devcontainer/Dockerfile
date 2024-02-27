FROM golang:1.22

WORKDIR /workspaces/${LOCAL_WORKSPACE_FOLDER_BASENAME}

RUN export DEBIAN_FRONTEND=noninteractive && apt update && apt -y upgrade

RUN export DEBIAN_FRONTEND=noninteractive && apt update && apt -y upgrade
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive && apt -y install --no-install-recommends \
    postgresql-client \
    pspg \
    tmux

RUN <<EOF cat >> ~/.bashrc

# Shorter prompt (fitting narrow VSCode windows without wrapping)
export PS1="\[\]finances:\w\$ \[\]"

# Default PostgreSQL Connection
# Will automatically connect to the devcontainer's attached PostgreSQL
# database (see "docker-compose.yml") by typing "psql".
export PGHOST=127.0.0.1
export PGPORT=40002
export PGUSER=postgres
export PGDB=finances_dev

# POSTGRESQL_URL is used as default for the `migrate` CLI.
export POSTGRESQL_URL='postgres://$PGUSER:@$PGHOST:$PGPORT/$PGDB?sslmode=disable'

export PAGER=pspg
EOF

# Install tools needed by the Go VSCODE plugin
RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v github.com/cweill/gotests/gotests@v1.6.0
RUN go install -v github.com/fatih/gomodifytags@v1.16.0
RUN go install -v github.com/josharian/impl@v1.1.0
RUN go install -v github.com/haya14busa/goplay/cmd/goplay@v1.0.0
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest
RUN go install -v honnef.co/go/tools/cmd/staticcheck@latest

# Install tools needed by the project's development operations
RUN go install -v github.com/a-h/templ/cmd/templ@latest
RUN go install -v github.com/DarthSim/overmind@latest
RUN go install -v github.com/cosmtrek/air@v1.50.0
RUN go install -v honnef.co/go/tools/cmd/staticcheck@latest
RUN go install -v -tags='no_clickhouse no_duckdb no_mssql no_mysql no_sqlite3 no_libsql no_vertica no_ydb' github.com/pressly/goose/v3/cmd/goose@latest
