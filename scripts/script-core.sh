#!/usr/bin/env bash
cd src/github.com/danielsussa/stock-listener

export TCP_URL="127.0.0.1:8081"

CompileDaemon -command="build/api-core" -directory="api-core" -build="go build -a -installsuffix cgo -o ../build/api-core cmd/main.go"

exec "$@"