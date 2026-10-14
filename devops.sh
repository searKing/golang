#! /bin/bash
#
# Copyright 2024 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

# Release modules by module sets defined in versions.yaml, as opentelemetry-go does.
# See https://github.com/open-telemetry/opentelemetry-go/blob/main/RELEASING.md

set -euo pipefail

cur_path=$(cd "$(dirname "$0")";pwd)
cd "${cur_path}"

export GOWORK=off

usage() {
  cat <<EOF
Usage: $0 COMMAND [OPTIONS]

Commands:
  tidy [-d DIR]                  run 'go mod tidy' in every module under DIR (default .)
  crosslink                      add replace directives of intra-repo requires and prune unused ones
  verify                         verify versions.yaml lists every module of the repo
  prerelease [-m MODSET]         update intra-repo requires to the versions in versions.yaml,
                                 run 'go mod tidy' and commit to a new branch prerelease_<MODSET>_<VERSION>
  tag [-m MODSET] [-c COMMIT]    tag modules of MODSET at COMMIT (default HEAD) by the version in versions.yaml
  untag [-m MODSET]              delete local tags of MODSET with the version in versions.yaml
  push [-m MODSET]               push tags of MODSET with the version in versions.yaml to origin

MODSET defaults to stable-v1.

Release steps:
  1. Bump the version of MODSET in versions.yaml, and commit it.
  2. $0 prerelease -m MODSET, then merge branch prerelease_<MODSET>_<VERSION>.
     Modules require each other by replace directives, run '$0 crosslink' if a new one is required.
  3. $0 tag -m MODSET -c COMMIT, where COMMIT is the merged commit.
  4. $0 push -m MODSET.
EOF
}

multimod() {
  go tool -modfile="${cur_path}/internal/tools/go.mod" multimod "$@"
}

crosslink() {
  go tool -modfile="${cur_path}/internal/tools/go.mod" crosslink "$@"
}

# modset_tags prints tags of the module set $1 with its version in versions.yaml.
modset_tags() {
  local modset=$1
  local version
  version=$(awk -v set="  ${modset}:" '$0 == set {found=1; next} found && $1 == "version:" {print $2; exit}' versions.yaml)
  if [[ -z "${version}" ]]; then
    echo "module set ${modset} not found in versions.yaml" >&2
    return 1
  fi
  awk -v set="  ${modset}:" '
    $0 == set {found=1; next}
    found && /^  [^ ]/ {exit}
    found && $1 == "-" {print $2}
  ' versions.yaml | while read -r mod; do
    if [[ "${mod}" == "github.com/searKing/golang" ]]; then
      echo "${version}"
    else
      echo "${mod#github.com/searKing/golang/}/${version}"
    fi
  done
}

cmd=${1:-}
if [[ -z "${cmd}" || "${cmd}" == "-h" || "${cmd}" == "help" ]]; then
  usage
  exit 0
fi
shift

args_dir="."
args_modset="stable-v1"
args_commit="HEAD"
while getopts d:m:c:h option; do
  case $option in
  d) args_dir=$OPTARG ;;
  m) args_modset=$OPTARG ;;
  c) args_commit=$OPTARG ;;
  h) usage; exit 0 ;;
  ?) usage; exit 1 ;;
  esac
done

case "${cmd}" in
tidy)
  find "${args_dir}" -type f -name go.mod -not -path '*/.*' -not -path '*/testdata/*' -print0 |
    while IFS= read -r -d '' mod; do
      dir=$(dirname "${mod}")
      echo "${dir#./}"
      (cd "${dir}" && go mod tidy)
    done
  echo "go mod tidy finished."
  ;;
crosslink)
  crosslink --root="${cur_path}" --prune
  ;;
verify)
  multimod verify
  ;;
prerelease)
  # multimod requires a clean working tree without untracked files, so prerelease in a clone of HEAD.
  if ! git diff --quiet HEAD; then
    echo "commit or stash changes of tracked files first" >&2
    exit 1
  fi
  head=$(git rev-parse HEAD)
  tmp=$(mktemp -d)
  trap 'rm -rf "${tmp}"' EXIT
  git clone -q --no-hardlinks "${cur_path}" "${tmp}/repo"
  (
    cd "${tmp}/repo"
    git checkout -q -B release_base "${head}"
    multimod verify
    multimod prerelease -m "${args_modset}"
  )
  git fetch -q "${tmp}/repo" 'refs/heads/prerelease_*:refs/heads/prerelease_*'
  git branch --list 'prerelease_*' --contains "${head}"
  ;;
tag)
  # Lightweight tags as the released ones, for 'multimod tag' signs tags by gpg always.
  commit=$(git rev-parse --verify "${args_commit}^{commit}")
  tags=()
  while read -r tag; do
    if git rev-parse -q --verify "refs/tags/${tag}" >/dev/null; then
      echo "tag ${tag} already exists" >&2
      exit 1
    fi
    tags+=("${tag}")
  done < <(modset_tags "${args_modset}")
  for tag in "${tags[@]}"; do
    git tag "${tag}" "${commit}"
    echo "${tag}"
  done
  echo "${#tags[@]} tags created at ${commit}. Run '$0 push -m ${args_modset}' to push them."
  ;;
untag)
  modset_tags "${args_modset}" | while read -r tag; do
    if git rev-parse -q --verify "refs/tags/${tag}" >/dev/null; then
      git tag -d "${tag}"
    fi
  done
  ;;
push)
  tags=()
  while read -r tag; do
    tags+=("refs/tags/${tag}")
  done < <(modset_tags "${args_modset}")
  printf '%s\n' "${tags[@]#refs/tags/}"
  read -r -p "Push the above ${#tags[@]} tags to origin? Enter y to continue: " reply
  if [[ "${reply}" == "y" ]]; then
    git push origin "${tags[@]}"
  else
    echo "Tag push cancelled."
  fi
  ;;
*)
  usage
  exit 1
  ;;
esac
