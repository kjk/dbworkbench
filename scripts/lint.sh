#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

echo "running eslint"
./node_modules/.bin/eslint jsx/*.js* jsx/reactable/paginator.jsx

echo "running go vet"
godep go vet github.com/kjk/dbworkbench
