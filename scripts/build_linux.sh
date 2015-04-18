#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

#godep go tool vet .

BINDATA_IGNORE=$(git ls-files -io --exclude-standard static/... | sed 's/^/-ignore=/;s/[.]/[.]/g')
go-bindata $BINDATA_IGNORE -ignore=[.]gitignore -ignore=[.]gitkeep static/...
GOOS=linux GOARCH=amd64 godep go build -o dbworkbench_linux
