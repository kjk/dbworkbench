#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

#echo "go vet"
#go tool vet -printfuncs=LogInfof,LogErrorf,LogVerbosef .
echo "go build"
godep go build -o quicknotes
#go build -o quicknotes
#gdep go build -race -o quicknotes

BINDATA_IGNORE=$(git ls-files -io --exclude-standard static/... | sed 's/^/-ignore=/;s/[.]/[.]/g')
go-bindata $BINDATA_IGNORE -ignore=[.]gitignore -ignore=[.]gitkeep static/...
godep go build -o dbworkbench
echo "starging dbworkbench"
./dbworkbench || true
rm dbworkbench
