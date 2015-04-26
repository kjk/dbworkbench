#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

#echo "go vet"
#go tool vet -printfuncs=LogInfof,LogErrorf,LogVerbosef .
echo "building"
godep go build -o dbworkbench
#gdep go build -race -o dbworkbench

echo "starting dbworkbench"
./dbworkbench --local || true
rm dbworkbench
