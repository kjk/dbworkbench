#!/bin/bash
set -u -e -o pipefail

cd website
go vet ./...

go test .

go build
./website "$@"|| true
rm website
