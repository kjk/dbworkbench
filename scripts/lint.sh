#!/bin/bash
set -u -e -o pipefail

echo "running eslint"
./node_modules/.bin/eslint js/*.js*

echo "running go vet"
godep go vet github.com/kjk/dbworkbench
