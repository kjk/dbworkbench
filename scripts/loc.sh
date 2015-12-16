#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

rm resources.go
wc -l s/*.html jsx/*.js* s/css/main.css
echo && wc -l *.go
