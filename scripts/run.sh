#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

#echo "go vet"
#go tool vet -printfuncs=LogInfof,LogErrorf,LogVerbosef .
echo "go build"
godep go build -o dbworkbench
#gdep go build -race -o dbworkbench

echo "starging dbworkbench"
./dbworkbench || true
rm dbworkbench
