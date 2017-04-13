#!/bin/bash
set -u -e -o pipefail

echo "formating js code with prettier"
./node_modules/.bin/prettier --trailing-comma es5 --write js/*js* js/tests/*.js js/alert/*.js
