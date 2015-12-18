#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

godep go vet github.com/kjk/dbworkbench

#TODO: use go tool vet so that I can pass printfuncs, but needs
#to filter out Godeps becase . is recursive
#godep go tool vet -printfuncs=LogInfof,LogErrorf,LogVerbosef .

rm -rf dbworkbench.test
godep go test -c ./...
rm -rf dbworkbench.test

godep go build -o dbherohelper
rm -rf dbherohelper
