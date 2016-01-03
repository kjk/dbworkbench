#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

./node_modules/.bin/gulp tests

#./node_modules/.bin/mocha jsx/tests/*.js
