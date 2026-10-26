#!/bin/bash

set -e

DATE=$(date +%Y%m%d.%H%M%S)

printf '> Setting up:\n'
if [ ! -f go.mod ]; then
	go mod init
fi
go mod tidy
go mod vendor
printf "> Building app '%s':\n" "$APP"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w -extldflags '-static' -X main.BuildTag=$DATE" -a -o /bin/$APP .
printf "> Done!\n"
