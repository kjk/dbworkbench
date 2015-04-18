#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

#godep go tool vet .

GOOS=linux GOARCH=amd64 godep go build -o dbworkbench_linux
