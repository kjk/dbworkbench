#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

./node_modules/.bin/gulp build_and_watch
