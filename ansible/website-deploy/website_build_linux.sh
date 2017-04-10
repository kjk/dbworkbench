#!/bin/bash
set -u -e -o pipefail

cd website
go tool vet .

GOOS=linux GOARCH=amd64 go build -o website_linux
