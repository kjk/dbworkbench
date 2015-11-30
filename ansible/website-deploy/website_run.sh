#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

cd /home/dbworkbench/www/app/current
exec ./website -s "$@" &>>crash.log
