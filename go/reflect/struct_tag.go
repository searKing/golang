package reflect

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
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
type StructTag map[string]SubStructTag

// ParseStructTag parses a single struct field tag and returns the set of subTags.
// This code is based on the validateStructTag code in package
func ParseStructTag(tag string) (StructTag, error) {
	// This code is based on the validateStructTag code in package
	// cmd/vendor/golang.org/x/tools/go/analysis/passes/structtag/structtag.go.

	var tags = map[string]SubStructTag{}

	// Code borrowed from cmd/vendor/golang.org/x/tools/go/analysis/passes/structtag/structtag.go
	// Code borrowed from reflect/type.go
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
		}

		if strings.IndexByte(value, ' ') >= 0 {
			return nil, errTagValueSpace
		}

		res := strings.Split(value, ",")
		name := res[0]
		options := res[1:]
		if len(options) == 0 {
			options = nil
		}

		if _, has := tags[key]; has {
			return nil, errTagDuplicatedKey
		}
		tags[key] = SubStructTag{
			Key:     key,
			Name:    name,
			Options: options,
		}
	}

	return tags, nil
}

// Get returns the tag associated with the given key. If the key is present
// in the tag the value (which may be empty) is returned. Otherwise the
// returned value will be the empty string. The ok return value reports whether
// the tag exists or not (which the return value is nil).
func (t StructTag) Get(key string) (tag SubStructTag, ok bool) {
	tag, ok = t[key]
	return
}

// Set sets the given tag. If the tag key already exists it'll override it
func (t StructTag) Set(subTag SubStructTag) error {
	if subTag.Key == "" {
		return errKeyNotSet
	}

	t[subTag.Key] = subTag
	return nil
}

// AddOptions adds the given option for the given key.
// It appends to any existing options associated with key.
func (t StructTag) AddOptions(key string, options ...string) {
	val, _ := t[key]
	val.AddOptions(options...)
	t[key] = val
}

// DeleteOptions deletes the given options for the given key
func (t StructTag) DeleteOptions(key string, options ...string) {

	val, has := t[key]
	if !has {
		return
	}
	val.DeleteOptions(options...)
	t[key] = val
}

// Delete deletes the tag for the given keys
func (t StructTag) Delete(keys ...string) {
	for _, key := range keys {
		delete(t, key)
	}
}

// StructTag returns a slice of subTags.
func (t StructTag) Tags() []SubStructTag {
	var tags []SubStructTag
	for _, tag := range t {
		tags = append(tags, tag)
	}
	return tags
}

// StructTag returns a slice of subTags.
func (t StructTag) Keys() []string {
	var keys []string
	for key, _ := range t {
		keys = append(keys, key)
	}
	return keys
}

// String reassembles the subTags into a valid literal tag field representation
func (t *StructTag) String() string {
	tags := t.Tags()
	if len(tags) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for i, tag := range t.Tags() {
		buf.WriteString(tag.String())
		if i != len(tags)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

var checkTagDups = []string{"json", "xml"}
var checkTagSpaces = map[string]bool{"json": true, "xml": true, "asn1": true}
