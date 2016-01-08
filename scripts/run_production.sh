#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

echo "running jsfmt"
./node_modules/.bin/esformatter -i js/*js* *.js

. scripts/lint.sh

echo "runnig gulp"
./node_modules/.bin/gulp prod

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

echo "starting dbherohelper in production mode"
./dbherohelper || true
rm dbherohelper
