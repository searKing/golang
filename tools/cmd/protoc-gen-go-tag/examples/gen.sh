#
# Copyright 2020 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

protoc -I . -I ../../../../../../../ --go-tag_out=paths=source_relative:. *.proto
