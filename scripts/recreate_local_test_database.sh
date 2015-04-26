#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

psql postgres <data/booktown.sql
echo "database is: postgres://localhost/booktown"
