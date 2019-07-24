#!/usr/bin/env bash

git clone https://github.com/searKing/golang.git .travis.workspace/golang
pushd golang/
git filter-branch --prune-empty --subdirectory-filter tools/cmd/go-atomicvalue/ master
git remote set-url origin https://github.com/searKing/travis-ci.git
git push -f origin master:go-atomicvalue
popd
rm -Rf .travis.workspace/
