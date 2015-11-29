#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

./node_modules/.bin/gulp default
go run scripts/build_release.go
godep go vet github.com/kjk/dbworkbench

#TODO: use go tool vet so that I can pass printfuncs, but needs
#to filter out Godeps becase . is recursive
#godep go tool vet -printfuncs=LogInfof,LogErrorf,LogVerbosef .

#rm -rf dbworkbench.test
#godep go test ./...
#rm -rf dbworkbench.test

echo "building"
godep go build -o dbworkbench
#gdep go build -race -o dbworkbench

echo "starting dbworkbench in no dev mode"
./dbworkbench || true
rm dbworkbench
