#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

#godep go tool vet .

cd cmd/import_stack_overflow
GOOS=linux GOARCH=amd64 godep go build -o import_stack_overflow_linux
cd ../..
mv cmd/import_stack_overflow/import_stack_overflow_linux ansible
