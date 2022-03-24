// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"
	render_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin/render"
)

// DefaultRender returns the appropriate Binding instance based on the HTTP method
// and the content type.
func DefaultRender(ctx *gin.Context, obj interface{}) render.Render {
	switch ctx.ContentType() {
	case binding.MIMEJSON:
		return render_.JSONPB{Data: obj} // support proto3 if enabled
	case binding.MIMEXML, binding.MIMEXML2:
		return render.XML{Data: obj}
	case binding.MIMEPROTOBUF:
		return render.ProtoBuf{Data: obj}
	case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
		return render.MsgPack{Data: obj}
	case binding.MIMEYAML:
		return render.YAML{Data: obj}
	case binding.MIMEMultipartPOSTForm, binding.MIMEPOSTForm:
		fallthrough
	default:
		return render.String{Format: "%v", Data: []interface{}{obj}}
	}
}
