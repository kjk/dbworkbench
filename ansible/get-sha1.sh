#!/bin/bash
set -u -e -o pipefail

cd $GOPATH/src/github.com/kjk/dbworkbench
git log -1 --pretty=format:%H
