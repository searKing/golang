// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/searKing/golang/third_party/github.com/spf13/viper/proto/internal/protobuild"
)

// Merge merges src proto into dst map, which must be a message with the same descriptor.
//
// Populated scalar fields in src are copied to dst, while populated
// singular messages in src are merged into dst by recursively calling Merge.
// The elements of every list field in src is appended to the corresponded
// list fields in dst. The entries of every map field in src is copied into
// the corresponding map field in dst, possibly replacing existing entries.
// The unknown fields of src are appended to the unknown fields of dst.
func Merge(dst map[any]any, src proto.Message, opts ...MergeOption) {
	var opt mergeOptions
	opt.ApplyOptions(opts...)
	opt.mergeMessage(dst, src.ProtoReflect())
}

//go:generate go-option -type=merge
type merge = mergeOptions

// mergeOptions provides a namespace for merge functions, and can be
// exported in the future if we add user-visible merge options.
type mergeOptions struct {
	// UseJsonNames uses lowerCamelCase name in JSON field names of proto field name instead.
	UseJsonNames bool

	// UseEnumNumbers emits enum values as numbers.
	UseEnumNumbers bool

	// EmitUnknown specifies whether to emit unknown fields in the output.
	// If specified, the unmarshaler may be unable to parse the output.
	// The default is to exclude unknown fields.
	EmitUnknown bool
}

func (o mergeOptions) mergeMessage(dst map[any]any, src protoreflect.Message) {
	if dst == nil {
		return
	}
	src.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch {
		case fd.IsList():
			var m []any
			o.mergeList(&m, v.List(), fd)
			dst[o.fieldName(fd)] = o.appendSlice(dst[o.fieldName(fd)], m...)
		case fd.IsMap():
			m := map[any]any{}
			o.mergeMap(m, v.Map(), fd)
			dst[o.fieldName(fd)] = o.appendMap(dst[o.fieldName(fd)], m)
		case fd.Message() != nil:
			m := map[any]any{}
			o.mergeMessage(m, v.Message())
			dst[o.fieldName(fd)] = o.appendMessage(dst[o.fieldName(fd)], m)
		default:
			var m any
			o.mergeSingular(&m, v, fd)
			dst[o.fieldName(fd)] = m
		}
		o.clearOtherOneofFields(dst, fd)
		return true
	})

	// Merge unknown fields.
	if o.EmitUnknown {
		var v []byte
		o.mergeUnknown(&v, src.GetUnknown())
		dst[protobuild.Unknown] = o.appendUnknown(dst[protobuild.Unknown], v...)
	}
}

func (o mergeOptions) mergeList(dst *[]any, src protoreflect.List, fd protoreflect.FieldDescriptor) {
	// Merge semantics appends to the end of the existing list.
	for i, n := 0, src.Len(); i < n; i++ {
		switch v := src.Get(i); {
		case fd.Message() != nil:
			dstv := map[any]any{}
			o.mergeMessage(dstv, v.Message())
			*dst = append(*dst, dstv)
		default:
			var m any
			o.mergeSingular(&m, v, fd)
			*dst = append(*dst, m)
		}
	}
}

func (o mergeOptions) mergeMap(dst map[any]any, src protoreflect.Map, fd protoreflect.FieldDescriptor) {
	// Merge semantics replaces, rather than merges into existing entries.
	src.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		var kn, vn any
		o.mergeSingular(&kn, protoreflect.Value(k), fd.MapKey())
		o.mergeSingular(&vn, v, fd.MapValue())
		dst[kn] = vn
		return true
	})
}

// mergeSingular merges the given non-repeated field value. This includes
// all scalar types, enums, messages, and groups.
func (o mergeOptions) mergeSingular(dst *any, val protoreflect.Value, fd protoreflect.FieldDescriptor) {
	if !val.IsValid() {
		return
	}

	switch kind := fd.Kind(); kind {
	case protoreflect.BoolKind:
		*dst = val.Bool()
		return
	case protoreflect.StringKind:
		*dst = val.String()
		return
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		*dst = int32(val.Int())
		return
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		*dst = uint32(val.Uint())
		return
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		*dst = val.Int()
		return
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		*dst = val.Uint()
		return
	case protoreflect.FloatKind:
		*dst = float32(val.Float())
		return
	case protoreflect.DoubleKind:
		*dst = val.Float()
		return
	case protoreflect.BytesKind:
		*dst = o.cloneBytes(val).Bytes()
		return
	case protoreflect.EnumKind:
		desc := fd.Enum().Values().ByNumber(val.Enum())
		if o.UseEnumNumbers || desc == nil {
			*dst = int32(val.Enum())
		} else {
			*dst = string(desc.Name())
		}
		return
	case protoreflect.MessageKind, protoreflect.GroupKind:
		m := map[any]any{}
		o.mergeMessage(m, val.Message())
		*dst = o.appendMessage(*dst, m)
	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	}
}

// mergeUnknown parses the given []byte and merges fields out.
// This function assumes no encoding in the given []byte.
func (o mergeOptions) mergeUnknown(dst *[]byte, b []byte) {
	*dst = b
}

func (o mergeOptions) clearOtherOneofFields(m map[any]any, fd protoreflect.FieldDescriptor) {
	od := fd.ContainingOneof()
	if od == nil {
		return
	}
	num := fd.Number()
	for i := 0; i < od.Fields().Len(); i++ {
		if n := od.Fields().Get(i).Number(); n != num {
			delete(m, o.fieldName(od.Fields().Get(i)))
		}
	}
}

func (o mergeOptions) cloneBytes(v protoreflect.Value) protoreflect.Value {
	return protoreflect.ValueOfBytes(append([]byte{}, v.Bytes()...))
}

func (o mergeOptions) fieldName(fd protoreflect.FieldDescriptor) string {
	name := fd.TextName()
	if o.UseJsonNames {
		name = fd.JSONName()
	}
	return name
}

func (o mergeOptions) appendSlice(slice any, elems ...any) []any {
	if slice == nil {
		return elems
	}
	typ := reflect.TypeOf(slice)
	val := reflect.ValueOf(slice)
	if typ.Kind() == reflect.Slice {
		var s = make([]any, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			s = append(s, val.Index(i).Interface())
		}
		s = append(s, elems...)
		return s
	}
	return elems
}

func (o mergeOptions) appendMessage(m any, elems map[any]any) map[any]any {
	return o.appendMapOrMessage(m, elems, true)
}

func (o mergeOptions) appendMap(m any, elems map[any]any) map[any]any {
	return o.appendMapOrMessage(m, elems, false)
}

func (o mergeOptions) appendMapOrMessage(m any, elems map[any]any, message bool) map[any]any {
	if m == nil {
		return elems
	}
	typ := reflect.TypeOf(m)
	val := reflect.ValueOf(m)
	if typ.Kind() == reflect.Map {
		var s = make(map[any]any, val.Len())
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			s[k.Interface()] = v.Interface()
		}
		for k, v := range elems {
			if message {
				if vv, ok := v.(map[any]any); ok {
					v = o.appendMapOrMessage(s[k], vv, message)
				}
			}
			s[k] = v
		}
		return s
	}
	return elems
}

func (o mergeOptions) appendUnknown(slice any, elems ...byte) []byte {
	if slice == nil {
		return elems
	}
	if b, ok := slice.(protoreflect.RawFields); ok {
		b = append(b, elems...)
		return b
	}
	if b, ok := slice.([]byte); ok {
		b = append(b, elems...)
		return b
	}
	return elems
}
