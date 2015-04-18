#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

cd /home/dbworkbench/www/app/current
exec ./dbworkbench "$@" &>>crash.log
