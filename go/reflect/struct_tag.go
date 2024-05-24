// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"bytes"
	"errors"
	"sort"
	"strconv"
	"strings"

	strings_ "github.com/searKing/golang/go/strings"
)

// cmd/vendor/golang.org/x/tools/go/analysis/passes/structtag/structtag.go
var (
	errTagSyntax        = errors.New("bad syntax for struct tag pair")
	errTagKeySyntax     = errors.New("bad syntax for struct tag key")
	errTagValueSyntax   = errors.New("bad syntax for struct tag value")
	errTagValueSpace    = errors.New("suspicious space in struct tag value")
	errTagSpace         = errors.New("key:\"value\" pairs not separated by spaces")
	errTagDuplicatedKey = errors.New("duplicated for struct tag key")

	errKeyNotSet   = errors.New("tag key does not exist")
	errTagNotExist = errors.New("tag does not exist")
)

// A StructTag is the tag string in a struct field, as reflect.StructTag .
//
// By convention, tag strings are a concatenation of
// optionally space-separated key:"value" pairs.
// Each key is a non-empty string consisting of non-control
// characters other than space (U+0020 ' '), quote (U+0022 '"'),
// and colon (U+003A ':').  Each value is quoted using U+0022 '"'
// characters and Go string literal syntax.
type StructTag struct {
	tagsByKey   map[string]SubStructTag
	orderedKeys []string
}

// ParseAstStructTag parses a single struct field tag of AST and returns the set of subTags.
// This code is based on the validateStructTag code in package
// tag `json:"name,omitempty"`, field.Tag.Value returned by AST
func ParseAstStructTag(tag string) (*StructTag, error) {
	if tag != "" {
		var err error
		tag, err = strconv.Unquote(tag)
		if err != nil {
			return nil, err
		}
	}
	return ParseStructTag(tag)
}

// ParseStructTag parses a single struct field tag and returns the set of subTags.
// This code is based on the validateStructTag code in package
// tag json:"name,omitempty", reflect.StructField.Tag returned by reflect
func ParseStructTag(tag string) (*StructTag, error) {
	// This code is based on the validateStructTag code in package
	// golang.org/x/tools/go/analysis/passes/structtag/structtag.go.
	var st = StructTag{
		tagsByKey: map[string]SubStructTag{},
	}

	// Code borrowed from golang.org/x/tools/go/analysis/passes/structtag/structtag.go
	// Code borrowed from reflect/type.go

	// This code is based on the StructTag.Get code in package reflect.
	n := 0
	for ; tag != ""; n++ {
		if n > 0 && tag != "" && tag[0] != ' ' {
			// More restrictive than reflect, but catches likely mistakes
			// like `x:"foo",y:"bar"`, which parses as `x:"foo" ,y:"bar"` with second key ",y".
			return nil, errTagSpace
		}
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 {
			return nil, errTagKeySyntax
		}
		if i+1 >= len(tag) || tag[i] != ':' {
			return nil, errTagSyntax
		}
		if tag[i+1] != '"' {
			return nil, errTagValueSyntax
		}
		key := tag[:i]
		if !IsValidTagKey(key) {
			return nil, errTagKeySyntax
		}
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			return nil, errTagValueSyntax
		}
		qvalue := tag[:i+1]
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return nil, errTagValueSyntax
		}

		if checkTagSpaces[key] {
			switch key {
			case "xml":
				// If the first or last character in the XML tag is a space, it is
				// suspicious.
				if strings.Trim(value, " ") != value {
					return nil, errTagValueSpace
				}

				// If there are multiple spaces, they are suspicious.
				if strings.Count(value, " ") > 1 {
					return nil, errTagValueSpace
				}

				// If there is no comma, skip the rest of the checks.
				comma := strings.IndexRune(value, ',')
				if comma < 0 {
					break
				}

				// If the character before a comma is a space, this is suspicious.
				if comma > 0 && value[comma-1] == ' ' {
					return nil, errTagValueSpace
				}
				//value = value[comma+1:]
			case "json":
				// JSON allows using spaces in the name, so skip it.
				comma := strings.IndexRune(value, ',')
				if comma < 0 {
					break
				}

				//value = value[comma+1:]
			}

			// If spaces exists in tag's value, it is suspicious.
			if strings.IndexByte(value, ' ') >= 0 {
				return nil, errTagValueSpace
			}
		}

		res := strings.Split(value, ",")
		name := res[0]
		options := res[1:]
		if len(options) == 0 {
			options = nil
		}

		if _, has := st.tagsByKey[key]; has {
			return nil, errTagDuplicatedKey
		}
		st.orderedKeys = append(st.orderedKeys, key)
		st.tagsByKey[key] = SubStructTag{
			Key:     key,
			Name:    name,
			Options: options,
		}
	}

	return &st, nil
}

// Get returns the tag associated with the given key. If the key is present
// in the tag the value (which may be empty) is returned. Otherwise the
// returned value will be the empty string. The ok return value reports whether
// the tag exists or not (which the return value is nil).
func (t StructTag) Get(key string) (tag SubStructTag, ok bool) {
	if t.orderedKeys == nil {
		return SubStructTag{}, false
	}
	tag, ok = t.tagsByKey[key]
	return
}

func (t *StructTag) appendOrderedKeysIfNotPresent(key string) {
	if t.tagsByKey == nil {
		t.tagsByKey = make(map[string]SubStructTag)
	}
	if _, has := t.tagsByKey[key]; !has {
		t.orderedKeys = append(t.orderedKeys, key)
	}
}

// Set sets the given tag. If the tag key already exists it'll override it
func (t *StructTag) Set(subTag SubStructTag) error {
	if subTag.Key == "" {
		return errKeyNotSet
	}
	t.appendOrderedKeysIfNotPresent(subTag.Key)
	t.tagsByKey[subTag.Key] = subTag
	return nil
}

// SetName sets the given name for the given key.
func (t *StructTag) SetName(key string, name string) {
	t.appendOrderedKeysIfNotPresent(key)

	val, _ := t.tagsByKey[key]
	val.Key = key
	val.Name = name
	t.tagsByKey[key] = val
}

// AddOptions adds the given option for the given key.
// It appends to any existing options associated with key.
func (t *StructTag) AddOptions(key string, options ...string) {
	t.appendOrderedKeysIfNotPresent(key)

	val, _ := t.tagsByKey[key]
	val.Key = key
	val.AddOptions(options...)
	t.tagsByKey[key] = val
}

// DeleteOptions deletes the given options for the given key
func (t *StructTag) DeleteOptions(key string, options ...string) {
	if t.tagsByKey == nil {
		return
	}

	val, has := t.tagsByKey[key]
	if !has {
		return
	}
	val.DeleteOptions(options...)
	t.tagsByKey[key] = val
}

// Delete deletes the tag for the given keys
func (t *StructTag) Delete(keys ...string) {
	if t.tagsByKey == nil {
		return
	}

	for _, key := range keys {
		delete(t.tagsByKey, key)
		strings_.SliceTrim(t.orderedKeys, key)
	}
}

// Keys returns a slice of subTags.
func (t StructTag) Keys() []string {
	var keys []string
	for key := range t.tagsByKey {
		keys = append(keys, key)
	}
	return keys
}

// SortedKeys returns a slice of subTags sorted by keys in increasing order.
func (t StructTag) SortedKeys() []string {
	keys := t.Keys()
	sort.Strings(keys)
	return keys
}

// OrderKeys returns a slice of subTags with original order.
func (t StructTag) OrderKeys() []string {
	return t.orderedKeys
}

// SelectedTags returns a slice of subTags in keys order.
func (t StructTag) SelectedTags(keys ...string) []SubStructTag {
	if len(keys) == 0 {
		return nil
	}
	var tags []SubStructTag
	for _, key := range keys {
		tag, has := t.tagsByKey[key]
		if !has {
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}

// SelectString reassembles the subTags selected by keys into a valid literal tag field representation
// tag json:"name,omitempty", reflect.StructField.Tag returned by reflect
func (t StructTag) SelectString(keys ...string) string {
	tags := t.SelectedTags(keys...)
	if len(tags) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for i, tag := range tags {
		buf.WriteString(tag.String())
		if i != len(tags)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

// SelectAstString reassembles the subTags selected by keys into a valid literal tag field representation
// tag json:"name,omitempty", reflect.StructField.Tag returned by reflect
func (t StructTag) SelectAstString(keys ...string) string {
	tag := t.SelectString(keys...)
	if tag == "" {
		return tag
	}
	return "`" + tag + "`"
}

// Tags returns a slice of subTags with original order.
func (t StructTag) Tags() []SubStructTag {
	return t.SelectedTags(t.Keys()...)
}

// SortedTags returns a slice of subTags sorted by keys in increasing order.
func (t StructTag) SortedTags() []SubStructTag {
	return t.SelectedTags(t.SortedKeys()...)
}

// OrderedTags returns a slice of subTags with original order.
func (t StructTag) OrderedTags() []SubStructTag {
	return t.SelectedTags(t.OrderKeys()...)
}

// String reassembles the subTags into a valid literal tag field representation
// key is random.
// tag json:"name,omitempty", reflect.StructField.Tag returned by reflect
func (t StructTag) String() string {
	return t.SelectString(t.Keys()...)
}

// SortedString reassembles the subTags into a valid literal tag field representation
// key is sorted by keys in increasing order.
// tag json:"name,omitempty", reflect.StructField.Tag returned by reflect
func (t StructTag) SortedString() string {
	return t.SelectString(t.SortedKeys()...)
}

// OrderedString reassembles the subTags into a valid literal tag field representation
// key is in the original order.
// tag json:"name,omitempty", reflect.StructField.Tag returned by reflect
func (t StructTag) OrderedString() string {
	return t.SelectString(t.OrderKeys()...)
}

// AstString reassembles the subTags into a valid literal ast tag field representation
// key is random.
// tag `json:"name,omitempty"`, field.Tag.Value returned by AST
func (t StructTag) AstString() string {
	return t.SelectAstString(t.Keys()...)
}

// SortedAstString reassembles the subTags into a valid literal ast tag field representation
// key is sorted by keys in increasing order.
// tag `json:"name,omitempty"`, field.Tag.Value returned by AST
func (t StructTag) SortedAstString() string {
	return t.SelectAstString(t.SortedKeys()...)
}

// OrderedAstString reassembles the subTags into a valid literal ast tag field representation
// key is in the original order.
// tag `json:"name,omitempty"`, field.Tag.Value returned by AST
func (t StructTag) OrderedAstString() string {
	return t.SelectAstString(t.OrderKeys()...)
}

var checkTagSpaces = map[string]bool{"json": true, "xml": true, "asn1": true}
