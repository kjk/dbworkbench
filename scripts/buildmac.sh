#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

godep go vet github.com/kjk/dbworkbench

./node_modules/.bin/gulp default

go run tools/build/*.go -gen-resources

godep go build -tags embeded_resources -o dbworkbench

cp dbworkbench mac/dbworkbench.exe

xcodebuild -parallelizeTargets -project mac/dbworkbench.xcodeproj/

rm mac/build/Release/dbworkbenchmacapp.zip
go run tools/build/*.go
