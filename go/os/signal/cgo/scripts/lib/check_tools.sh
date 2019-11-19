#!/bin/bash

if [[ ! ( $?SCRIPTS_LIB_CHECK_TOOLS_SH_ ) ]]; then
  return
fi

export SCRIPTS_LIB_CHECK_TOOLS_SH_="check_tools.sh"

# 获取当前脚本的相对路径文件名称
THIS_BASH_FILE="${BASH_SOURCE-$0}"
readonly THIS_BASH_FILE
# 获取当前脚本的相对路径目录
THIS_BASH_FILE_REF_DIR=$(dirname "${THIS_BASH_FILE}")
readonly THIS_BASH_FILE_REF_DIR
# 获取当前脚本的绝对路径目录
THIS_BASH_FILE_ABS_DIR=$(
  cd "${THIS_BASH_FILE_REF_DIR}" || exit
  pwd
)
readonly THIS_BASH_FILE_ABS_DIR

# include
pushd "${THIS_BASH_FILE_ABS_DIR}" 1>/dev/null 2>&1 || exit
. ./log.sh
popd 1>/dev/null 2>&1 || exit

# Sanity check that the right tools are accessible.
function check_tools() {
  for tool in "$@"; do
    q=$(command -v "${tool}") || die "didn't find ${tool}"
    log::error "$tool: $q"
  done
}
