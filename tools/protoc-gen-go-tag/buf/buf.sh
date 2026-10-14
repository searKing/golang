#!/usr/bin/env bash
#
# Copyright 2026 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

# Runs the buf command on the BSR module of tag.proto, such as:
#
#   ./buf.sh build
#   ./buf.sh push --label v1.2.147
#
# The module is laid out in a temporary directory, with tag.proto exported by protoc-gen-go-tag -export_proto,
# so that it is imported by its path github.com/searKing/golang/tools/protoc-gen-go-tag/tag/tag.proto.
#
# BUF is the buf command, defaults to buf in PATH.

set -euo pipefail

if [[ $# -eq 0 ]]; then
  echo "usage: $0 <buf command> [flags]" >&2
  exit 2
fi

cd "$(dirname "$0")"

tmp=$(mktemp -d)
trap 'rm -rf "${tmp}"' EXIT

cp buf.yaml README.md "${tmp}"
(cd .. && GOWORK=off go run . -export_proto "${tmp}")

cmd=$1
shift
${BUF:-buf} "${cmd}" "${tmp}" "$@"
