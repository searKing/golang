#!/usr/bin/env bash
#
# Copyright 2026 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

# Regenerates examples/example.pb.go and testdata/protos.binpb by protoc, run by go generate.
# Run 'go test -update' then to update golden files of testdata/golden.
#
# PROTOC_INCLUDE is the include directory of google/protobuf/descriptor.proto,
# defaults to the include directory installed with protoc.

set -euo pipefail

cd "$(dirname "$0")"

protoc_include=${PROTOC_INCLUDE:-$(dirname "$(dirname "$(command -v protoc)")")/include}
if [[ ! -f "${protoc_include}/google/protobuf/descriptor.proto" ]]; then
  echo "google/protobuf/descriptor.proto not found in ${protoc_include}, set PROTOC_INCLUDE to its include directory" >&2
  exit 1
fi

tmp=$(mktemp -d)
trap 'rm -rf "${tmp}"' EXIT

# tag.proto is imported by its path in the repo github.com/searKing/golang.
mkdir -p "${tmp}/include/github.com/searKing"
ln -s "$(cd ../.. && pwd)" "${tmp}/include/github.com/searKing/golang"
includes=(-I "${tmp}/include" -I "${protoc_include}")

GOWORK=off go build -o "${tmp}/protoc-gen-go-tag" .
protoc "${includes[@]}" -I examples \
  --plugin=protoc-gen-go-tag="${tmp}/protoc-gen-go-tag" \
  --go-tag_out=paths=source_relative:examples \
  example.proto

protos=()
while IFS= read -r proto; do
  protos+=("${proto}")
done < <(cd testdata/proto && find . -name '*.proto' | sed 's|^\./||' | sort)
protoc "${includes[@]}" -I examples -I testdata/proto \
  --include_imports --include_source_info \
  --descriptor_set_out=testdata/protos.binpb \
  example.proto "${protos[@]}"
