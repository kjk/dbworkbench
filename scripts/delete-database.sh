#!/bin/bash

# TODO: when schema is stabilized, delete this script

set -o nounset
set -o errexit
set -o pipefail

read -r -p "Are you sure you want to delete dbworkbench database? [y/N] " response
case $response in
    [yY][eE][sS]|[yY])
      psql postgres -c "DROP DATABASE dbworkbench"
      ;;
    *)
      echo "Didn't delete the database"
      ;;
esac
