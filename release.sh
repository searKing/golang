#! /bin/bash
#
# Copyright 2026 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

set -euo pipefail

cur_path=$(cd "$(dirname "$0")";pwd)
cd "${cur_path}"

echo "$0" "$*"

VERSION=v1.2.145
export VERSION

find . -type f -name go.mod -not -path './.*' -not -path './*/testdata/*' -exec bash -c 'd=$(dirname "$1"); d="${d#./}"; [ "$d" = "." ] && echo "$VERSION" || echo "$d/$VERSION"' _ {} \;

read -r -p "Create the above tags? Enter y to continue: " REPLY
if [ "$REPLY" = "y" ]; then
  find . -type f -name go.mod -not -path './.*' -not -path './*/testdata/*' -exec bash -c 'd=$(dirname "$1"); d="${d#./}"; tag="$d/$VERSION"; [ "$d" = "." ] && tag="$VERSION"; git tag "$tag"' _ {} \;
else
  echo "Tag creation cancelled."
fi