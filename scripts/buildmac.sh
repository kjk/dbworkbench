#!/bin/bash
set -u -e -o pipefail

# flags:
#   -upload : upload to s3
#   -beta   : build beta version (different location in s3)

. scripts/buildmachelper.sh

rm -rf mac/build

echo "running xcode..."
xcodebuild -parallelizeTargets -project mac/dbHero.xcodeproj/

codesign --force --deep --verbose -s "Developer ID Application: Krzysztof Kowalczyk (2LGSCEWRR9)" -f "mac/build/Release/dbHero.app"

codesign --verify --verbose "mac/build/Release/dbHero.app"

go run tools/build/*.go $@

rm resources.go
