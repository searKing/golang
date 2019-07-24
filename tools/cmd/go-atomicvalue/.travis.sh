#!/usr/bin/env bash

rm -Rf .travis.workspace/ || exit -1
git clone https://github.com/searKing/golang.git .travis.workspace/golang || exit -1
pushd .travis.workspace/golang/ || exit -1
git filter-branch --prune-empty --subdirectory-filter tools/cmd/go-atomicvalue/ master || exit -1
git remote set-url origin https://github.com/searKing/travis-ci.git || exit -1
git push -f origin master:go-atomicvalue || exit -1
popd || exit -1
rm -Rf .travis.workspace/ || exit -1
