#! /bin/bash
#
# Copyright 2024 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

set -euo pipefail

cur_path=$(cd "$(dirname "$0")";pwd)
cd "${cur_path}"

echo "$0" "$*"

args_dir="."
args_tag=""
args_clean_only=OFF
while getopts cd:t:h option; do
  case $option in
  c) args_clean_only=ON ;;
  d) args_dir=$OPTARG ;;
  t) args_tag=$OPTARG ;;
  h)
    echo './devops.sh -d GO-MOD-TIDY-ROOT-DIR-PATH' ;
    echo './devops.sh -d GO-MOD-TIDY-ROOT-DIR-PATH' -t 'TAG-VERSION' ;;
  ?) exit 1 ;;
  esac
done

if [[ -z "$args_dir" ]]; then
    echo "root dir should not be empty"
    exit 1
fi

if [[ -n "$args_tag" ]]; then
  git_args=('')
  if [ "${args_clean_only}"x = "ON"x ]; then
    git_args=(-d)
  fi
  cmd="path=\$(dirname \"\${1#./}\");if [ \"\$path\" == \".\" ]; then git tag ${git_args[*]} \"\${TAG}\"; else git tag ${git_args[*]} \"\${path}/\${TAG}\"; fi;"
  TAG="${args_tag}" find "${args_dir}" -type f -name "go.mod" -not -path "./.*" -not -path "./*/testdata/*" -exec bash -c "$cmd" sh {} \;
  if [ "${args_clean_only}"x = "ON"x ]; then
    echo "'$args_tag' tag deleted finished."
  else
    echo "'$args_tag' tag create finished."
  fi
  exit 1
fi

find "${args_dir}" -type f -name "go.mod" -not -path "./.*" -exec bash -c 'cd $(dirname "$1"); echo $(dirname "${1#./}"); go mod tidy' sh {} \;
echo "go mod tidy finished."
