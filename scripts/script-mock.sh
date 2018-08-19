#!/usr/bin/env bash
cd src/github.com/danielsussa/stock-listener
CompileDaemon -command="build/api-mock" -directory="api-mock" -build="go build -a -installsuffix cgo -o ../build/api-mock ."

exec "$@"