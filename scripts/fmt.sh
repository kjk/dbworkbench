#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

./node_modules/.bin/esformatter -i js/*js* js/tests/*.js

