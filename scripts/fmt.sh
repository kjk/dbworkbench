#!/bin/bash
set -u -e -o pipefail

echo "running esformatter"
./node_modules/.bin/esformatter -i js/*js* js/tests/*.js
