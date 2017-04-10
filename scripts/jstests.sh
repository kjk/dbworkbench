#!/bin/bash
set -u -e -o pipefail

./node_modules/.bin/gulp tests

#./node_modules/.bin/mocha js/tests/*.js
