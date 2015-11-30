#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

cd website
go tool vet .

GOOS=linux GOARCH=amd64 go build -o website_linux
