#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

wc -l static/index.html static/js/app.js static/css/app.css
echo && wc -l *.go
