#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

godep go vet github.com/kjk/dbworkbench

./node_modules/.bin/gulp default

go run tools/build/*.go -gen-resources

#TODO: use go tool vet so that I can pass printfuncs, but needs
#to filter out Godeps becase . is recursive
#godep go tool vet -printfuncs=LogInfof,LogErrorf,LogVerbosef .

#rm -rf dbworkbench.test
#godep go test ./...
#rm -rf dbworkbench.test

echo "building"
godep go build -tags embeded_resources -o dbherohelper
#gdep go build -tags embeded_resources -race -o dbherohelper

echo "starting dbherohelper in no dev mode"
./dbherohelper || true
rm dbherohelper
