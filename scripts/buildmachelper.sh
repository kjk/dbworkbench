#!/bin/bash
set -u -e -o pipefail

rm -rf mac/dbherohelper.exe dbherohelper.zip

. scripts/lint.sh

echo "running gulp prod"
./node_modules/.bin/gulp prod

echo "generating resources .zip file..."
go run tools/build/*.go -gen-resources

echo "building dbherohelper.exe..."
godep go build -tags embeded_resources -o mac/dbherohelper.exe

rm -rf dbherohelper.zip
