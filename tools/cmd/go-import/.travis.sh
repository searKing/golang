#!/usr/bin/env bash
set -o pipefail
set -o errexit
set -o nounset
# set -o xtrace

# 获取输入参数
THIS_BASE_PARAM="$*"
# 获取当前脚本的相对路径文件名称
THIS_BASH_FILE="${BASH_SOURCE-$0}"
# 获取当前脚本的相对路径目录
THIS_BASH_FILE_REF_DIR=$(dirname "${THIS_BASH_FILE}")
# 获取当前脚本的绝对路径目录
THIS_BASH_FILE_ABS_DIR=$(
  cd "${THIS_BASH_FILE_REF_DIR}" || exit
  pwd
)
# 获取当前脚本的名称
THIS_BASH_FILE_BASE_NAME=$(basename "${THIS_BASH_FILE}")
# 获取当前脚本绝对路径
THIS_BASH_FILE_ABS_PATH="${THIS_BASH_FILE_ABS_DIR}/${THIS_BASH_FILE_BASE_NAME}"
# 备份当前路径
STACK_ABS_DIR=$(pwd)
# 临时文件
# Install the working tree in a tempdir.
tmpdir=$(mktemp -d -t .build.XXXXXX)
function cleanup() {
  printf "Cleaning up %s..." "${tmpdir}"
  [ -d "${tmpdir}" ] && rm -Rf "${tmpdir}"
  printf "\r\033[KCleaning done."
  printf "\r\033[K"
}
trap cleanup EXIT
[ -d "${tmpdir}" ] && rm -Rf "${tmpdir}"
mkdir -p "${tmpdir}"
# 路径隔离
cd "${THIS_BASH_FILE_ABS_DIR}" || exit
[ -d "${tmpdir}"/ ] && rm -Rf "${tmpdir}"/ || exit
git clone https://github.com/searKing/golang.git "${tmpdir}/golang" || exit
pushd "${tmpdir}/golang" 1>/dev/null 2>&1 || exit
git filter-branch --prune-empty --subdirectory-filter tools/cmd/go-import/ master || exit
# reset and clean .git
git reset --hard
git for-each-ref --format="%(refname)" refs/original | xargs -n 1 git update-ref -d || exit
git reflog expire --expire=now --all || exit
git gc --aggressive --prune=now || exit

git remote set-url origin https://github.com/searKing/travis-ci.git || exit
git push -f origin master:go-import || exit
popd 1>/dev/null 2>&1 || exit
