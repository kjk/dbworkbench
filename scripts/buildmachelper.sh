#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

rm -rf mac/dbherohelper.exe dbherohelper.zip

godep go vet github.com/kjk/dbworkbench

./node_modules/.bin/gulp default

echo "generating resources .zip file..."
go run tools/build/*.go -gen-resources

echo "building dbherohelper.exe..."
godep go build -tags embeded_resources -o mac/dbherohelper.exe

rm -rf dbherohelper.zip
