#!/bin/bash
set -u -e -o pipefail

psql postgres <data/world.sql
echo "database is: postgres://localhost/world"
