#!/usr/bin/env bash

# 获取输入参数
THIS_BASE_PARAM="$@"
# 获取当前脚本的相对路径文件名称
THIS_BASH_FILE="${BASH_SOURCE-$0}"
# 获取当前脚本的相对路径目录
THIS_BASH_FILE_REF_DIR=`dirname ${THIS_BASH_FILE}`
# 获取当前脚本的绝对路径目录
THIS_BASH_FILE_ABS_DIR=`cd ${THIS_BASH_FILE_REF_DIR}; pwd`
# 获取当前脚本的名称
THIS_BASH_FILE_BASE_NAME=`basename ${THIS_BASH_FILE}`
# 获取当前脚本绝对路径
THIS_BASH_FILE_ABS_PATH="${THIS_BASH_FILE_ABS_DIR}/${THIS_BASH_FILE_BASE_NAME}"
# 备份当前路径
STACK_ABS_DIR=`pwd`
# 路径隔离
cd "${THIS_BASH_FILE_ABS_DIR}"

rm -Rf .travis.workspace/ || exit -1
git clone https://github.com/searKing/golang.git .travis.workspace/golang || exit -1
pushd .travis.workspace/golang/ || exit -1
git filter-branch --prune-empty --subdirectory-filter tools/cmd/go-syncmap/ master || exit -1
git remote set-url origin https://github.com/searKing/travis-ci.git || exit -1
git push -f origin master:go-syncmap || exit -1
popd || exit -1
rm -Rf .travis.workspace/ || exit -1
