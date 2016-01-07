#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

echo "running jsfmt"
./node_modules/.bin/esformatter -i jsx/reactable/*.jsx jsx/*js* *.js

echo "running eslint"
./node_modules/.bin/eslint jsx/*.js* jsx/reactable/*.jsx

echo "running go vet"
godep go vet github.com/kjk/dbworkbench

echo "running gulp prod"
./node_modules/.bin/gulp prod

#TODO: use go tool vet so that I can pass printfuncs, but needs
#to filter out Godeps becase . is recursive
#godep go tool vet -printfuncs=LogInfof,LogErrorf,LogVerbosef .

#rm -rf dbworkbench.test
#godep go test -c ./...
#rm -rf dbworkbench.test

godep go build -o dbherohelper
rm -rf dbherohelper
