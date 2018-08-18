#!/usr/bin/env bash
cd src/github.com/danielsussa/stock-listener
mkdir /build

echo 1-starting compile daemon
ls
CompileDaemon -command="build/api-core" -directory="api-core" -build="go build -a -installsuffix cgo -o ../build/api-core cmd/main.go"

exec "$@"