#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

godep go vet github.com/kjk/dbworkbench

# TODO: Should I do tests?

./node_modules/.bin/gulp default

go run scripts/build_release.go

godep go build -o dbworkbench

cp dbworkbench mac/dbworkbench.exe
cp dbworkbench.dat mac/

go run tools/build/main.go tools/build/util.go tools/build/cmd.go tools/build/s3.go tools/build/win.go -no-clean-check

xcodebuild -parallelizeTargets -project mac/dbworkbench.xcodeproj/