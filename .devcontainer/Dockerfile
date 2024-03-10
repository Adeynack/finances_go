FROM golang:1.22

WORKDIR /workspaces/${LOCAL_WORKSPACE_FOLDER_BASENAME}

RUN export DEBIAN_FRONTEND=noninteractive && apt update && apt -y upgrade

RUN export DEBIAN_FRONTEND=noninteractive && apt update && apt -y upgrade
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive && apt -y install --no-install-recommends \
    postgresql-client \
    pspg \
    tmux

RUN <<EOF cat >> ~/.bashrc

export LC_ALL=C
export LANGUAGE=$LANG

# Shorter prompt (fitting narrow VSCode windows without wrapping)
export PS1="\[\]finances:\w\$ \[\]"

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

# NodeJS
RUN curl -fsSL https://deb.nodesource.com/setup_21.x | bash - && \
    apt install -y nodejs
