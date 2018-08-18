#!/usr/bin/env bash
cd src/github.com/danielsussa/stock-listener
mkdir /build

echo 1-starting compile daemon
ls
CompileDaemon -command="build/api-mock" -directory="api-mock" -build="go build -a -installsuffix cgo -o ../build/api-mock ."

exec "$@"