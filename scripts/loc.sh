#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

wc -l s/*.html jsx/*.js* s/css/app_react.css
echo && wc -l *.go
