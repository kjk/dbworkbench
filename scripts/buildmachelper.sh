#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

rm -rf mac/dbherohelper.exe dbherohelper.zip

echo "running go vet"
godep go vet github.com/kjk/dbworkbench

echo "running gulp prod"
./node_modules/.bin/gulp prod

echo "generating resources .zip file..."
go run tools/build/*.go -gen-resources

echo "building dbherohelper.exe..."
godep go build -tags embeded_resources -o mac/dbherohelper.exe

rm -rf dbherohelper.zip
