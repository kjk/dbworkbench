#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

echo "running jsfmt"
./node_modules/.bin/esformatter -i jsx/reactable/*.jsx jsx/*js* *.js

. scripts/lint.sh

#TODO: use go tool vet so that I can pass printfuncs, but needs
#to filter out Godeps becase . is recursive
#godep go tool vet -printfuncs=LogInfof,LogErrorf,LogVerbosef .

#rm -rf dbworkbench.test
#godep go test ./...
#rm -rf dbworkbench.test

echo "building"
godep go build -o dbherohelper
#gdep go build -race -o dbherohelper

echo "starting dbherohelper in dev mode"
./dbherohelper -dev "$@" || true
rm dbherohelper
