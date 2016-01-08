#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

rm -f resources.go || true

echo && wc -l scripts/*.*
echo && wc -l mac/dbHero/*.swift
echo && wc -l win/dbhero/*.cs
echo && wc -l s/*.html sass/main.scss
echo && wc -l js/*.js* #js/reactable/*.js*
echo && wc -l *.go website/*.go
