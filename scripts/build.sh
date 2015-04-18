#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

#godep go tool vet .

godep go build -o dbworkbench
#rm -rf dbworkbench
