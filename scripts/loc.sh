#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

wc -l s/index.html s/js/app.js s/css/app.css
echo && wc -l *.go
