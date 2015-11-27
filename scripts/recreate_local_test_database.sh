#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

psql postgres <data/world.sql
echo "database is: postgres://localhost/world"
