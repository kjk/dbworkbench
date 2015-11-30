#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

godep go vet github.com/kjk/dbworkbench

# TODO: Should I do tests?

./node_modules/.bin/gulp default

go run tools/build/*.go -gen-resources

godep go build -o dbworkbench

cp dbworkbench mac/dbworkbench.exe
cp dbworkbench.dat mac/

xcodebuild -parallelizeTargets -project mac/dbworkbench.xcodeproj/

go run tools/build/*.go