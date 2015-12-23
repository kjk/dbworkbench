#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

rm -f resources.go || true
<<<<<<< HEAD
echo && wc -l mac/dbHero/*.swift
echo && wc -l win/dbhero/*.cs
echo && wc -l s/*.html sass/main.scss
echo && wc -l jsx/*.js* jsx/lib/*.js* jsx/lib/reactable/*.js*
echo && wc -l *.go website/*.go
echo && wc -l scripts/*.*
