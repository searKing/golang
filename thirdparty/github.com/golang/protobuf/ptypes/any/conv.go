package any

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	_struct "github.com/searKing/golang/thirdparty/github.com/golang/protobuf/ptypes/struct"
)

// ToProtoAny converts v, which must marshal into a JSON object,
// into a Google Any proto.
func ToProtoAny(data interface{}) (*any.Any, error) {
	if data == nil {
		return &any.Any{}, nil
	}
	var datapb proto.Message
	switch data.(type) {
	case proto.Message:
		datapb = data.(proto.Message)
	default:
		dataStructpb, err := _struct.ToProtoStruct(data)
		if err != nil {
			return nil, err
		}
		datapb = dataStructpb
	}
	return ptypes.MarshalAny(datapb)
}
