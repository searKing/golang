package prettyjson

import (
	"bytes"
	"encoding"
	"encoding/base64"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"unicode/utf8"

	bytes_ "github.com/searKing/golang/go/bytes"
	strings_ "github.com/searKing/golang/go/strings"
)

// Marshal serializes v using encoding/json/v2 and applies the configured
// truncation rules.
func Marshal(v any, options ...EncOptsOption) ([]byte, error) {
	return jsonv2.Marshal(v, EncOptions(options...)...)
}

// MarshalIndent is like [Marshal] but applies [jsonv1.Indent] to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func MarshalIndent(v any, prefix, indent string, options ...EncOptsOption) ([]byte, error) {
	b, err := Marshal(v, options...)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := jsonv1.Indent(&buf, b, prefix, indent); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalWrite is like [Marshal] but writes the JSON to out.
// It does not terminate the output with a newline.
func MarshalWrite(out io.Writer, in any, options ...EncOptsOption) (err error) {
	return jsonv2.MarshalWrite(out, in, EncOptions(options...)...)
}

// MarshalEncode is like [Marshal] but encodes the JSON to out.
// The options provided take precedence over options already applied on out,
// such as indentation, and only apply for the duration of the call.
func MarshalEncode(out *jsontext.Encoder, in any, options ...EncOptsOption) (err error) {
	return jsonv2.MarshalEncode(out, in, EncOptions(options...)...)
}

// EncOptions returns the json/v2 options that [Marshal] applies for options,
// so that they can be passed to json/v2 functions directly.
func EncOptions(options ...EncOptsOption) []jsonv2.Options {
	var opts encOpts
	opts.ApplyOptions(options...)

	marshalers := truncateMarshalers(opts)

	var jsonOpts []jsonv2.Options
	if marshalers != nil {
		jsonOpts = append(jsonOpts, jsonv2.WithMarshalers(marshalers))
	}
	return append(jsonOpts, formatOptions(opts)...)
}

// formatOptions returns the json/v2 options of opts other than truncation.
func formatOptions(opts encOpts) []jsonv2.Options {
	jsonOpts := []jsonv2.Options{
		jsontext.EscapeForHTML(opts.escapeHTML),
		jsonv2.Deterministic(true),
	}
	if opts.omitEmpty {
		// encoding/json's v1 compatibility option is an encoding/json/v2
		// Options value and can therefore be passed directly to json/v2.
		// It only applies to fields tagged with omitempty.
		jsonOpts = append(jsonOpts, jsonv1.OmitEmptyWithLegacySemantics(true))
	}
	if opts.omitEmpty || opts.omitZero {
		jsonOpts = append(jsonOpts, jsonv2.OmitZeroStructFields(true))
	}
	return jsonOpts
}

// AdvancedTruncateMarshalers returns the type-specific marshalers that [Marshal]
// truncates values with for options, so that they can be joined with other
// marshalers by [jsonv2.JoinMarshalers]. Options other than truncation, such as
// escaping or omitting fields, are ignored; see [EncOptions] for those.
func AdvancedTruncateMarshalers(options ...EncOptsOption) *jsonv2.Marshalers {
	var opts encOpts
	opts.ApplyOptions(options...)
	return truncateMarshalers(opts)
}

// truncateMarshalers returns the type-specific marshalers for opts. string is
// handled by a dedicated marshaler, and all other types by the catch-all "any"
// marshaler, which uses reflect.Kind to distinguish maps from structs.
func truncateMarshalers(opts encOpts) *jsonv2.Marshalers {
	return jsonv2.JoinMarshalers(
		stringMarshalers(opts),
		anyMarshalers(opts),
	)
}

func anyMarshalers(opts encOpts) *jsonv2.Marshalers {
	if opts.truncateString <= 0 && opts.truncateBytes <= 0 && opts.truncateSliceOrArray <= 0 && opts.truncateMap <= 0 {
		return nil
	}
	truncateOutput := opts.truncateMarshalerOutput && (opts.truncateString > 0 || opts.truncateSliceOrArray > 0)

	return jsonv2.MarshalToFunc(func(enc *jsontext.Encoder, v any) error {
		// json/v2 always passes a non-nil pointer to the value. Nested pointers
		// and interfaces are left to json/v2, which calls back with their elements.
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Pointer || rv.IsNil() {
			return errors.ErrUnsupported
		}
		rv = rv.Elem()

		switch rv.Kind() {
		case reflect.String:
			// Named string types; string itself is left to stringMarshalers.
			if opts.truncateString <= 0 || isObjectName(enc) || hasMarshalMethod(rv.Type()) {
				return errors.ErrUnsupported
			}
			return enc.WriteToken(jsontext.String(truncateStringValue(rv.String(), opts)))
		case reflect.Slice, reflect.Array:
			if rv.Kind() == reflect.Slice && rv.IsNil() {
				return errors.ErrUnsupported
			}
			// Byte slices and arrays are base64-encoded by default. Only slices
			// of byte are truncated as bytes; byte arrays are left untouched.
			if rv.Type().Elem().Kind() == reflect.Uint8 {
				if rv.Kind() != reflect.Slice || rv.Type().Elem() != byteType ||
					opts.truncateBytes <= 0 || hasMarshalMethod(rv.Type()) {
					return errors.ErrUnsupported
				}
				return marshalTruncatedBytes(enc, rv.Bytes(), opts)
			}
		case reflect.Map:
			if rv.IsNil() {
				return errors.ErrUnsupported
			}
		default:
			// Structs and all other types use the default marshaler.
			return errors.ErrUnsupported
		}

		if truncateOutput && hasMarshalMethod(rv.Type()) {
			return marshalTruncatedMarshalerOutput(enc, v, opts)
		}

		n := rv.Len()
		if rv.Kind() == reflect.Map {
			limit, limitIfMoreThan := truncateTo(opts.truncateMap, opts.truncateMapIfMoreThan)
			if opts.truncateMap <= 0 || n <= limitIfMoreThan || n <= limit || hasMarshalMethod(rv.Type()) {
				return errors.ErrUnsupported
			}
			return marshalTruncatedMap(enc, rv, n, limit, opts)
		}
		limit, limitIfMoreThan := truncateTo(opts.truncateSliceOrArray, opts.truncateSliceOrArrayIfMoreThan)
		if opts.truncateSliceOrArray <= 0 || n <= limitIfMoreThan || n <= limit || hasMarshalMethod(rv.Type()) {
			return errors.ErrUnsupported
		}
		return marshalTruncatedSliceOrArray(enc, rv, n, limit, opts)
	})
}

var byteType = reflect.TypeFor[byte]()

// marshalerTypes are the method interfaces json/v2 prefers over its default
// representation of a type.
var marshalerTypes = []reflect.Type{
	reflect.TypeFor[jsonv2.MarshalerTo](),
	reflect.TypeFor[jsonv2.Marshaler](),
	reflect.TypeFor[encoding.TextAppender](),
	reflect.TypeFor[encoding.TextMarshaler](),
}

// hasMarshalMethod reports whether json/v2 marshals values of t through one of
// their methods. json/v2 values are addressable, so pointer receivers count.
func hasMarshalMethod(t reflect.Type) bool {
	pt := reflect.PointerTo(t)
	for _, it := range marshalerTypes {
		if t.Implements(it) || pt.Implements(it) {
			return true
		}
	}
	return false
}

// marshalTruncatedMarshalerOutput marshals v through its own method, then
// writes the output with its strings and arrays truncated. Objects keep all
// their members, as they cannot be told apart from structs.
func marshalTruncatedMarshalerOutput(enc *jsontext.Encoder, v any, opts encOpts) error {
	out, err := jsonv2.Marshal(v, formatOptions(opts)...)
	if err != nil {
		return err
	}
	return writeTruncatedValue(enc, out, opts)
}

// writeTruncatedValue writes the JSON value val, truncating the strings and
// arrays in it like Go strings and slices.
func writeTruncatedValue(enc *jsontext.Encoder, val jsontext.Value, opts encOpts) error {
	switch val.Kind() {
	case '"':
		s, err := jsontext.AppendUnquote(nil, val)
		if err != nil {
			return err
		}
		return enc.WriteToken(jsontext.String(truncateStringValue(string(s), opts)))
	case '[':
		return writeTruncatedArray(enc, val, opts)
	case '{':
		return writeTruncatedObject(enc, val, opts)
	default:
		return enc.WriteValue(val)
	}
}

// writeTruncatedArray writes the JSON array val, keeping as many elements as
// a Go slice of the same length keeps.
func writeTruncatedArray(enc *jsontext.Encoder, val jsontext.Value, opts encOpts) error {
	dec := jsontext.NewDecoder(bytes.NewReader(val))
	if _, err := dec.ReadToken(); err != nil {
		return err
	}
	var n int
	for dec.PeekKind() != ']' {
		if err := dec.SkipValue(); err != nil {
			return err
		}
		n++
	}

	limit := n
	if opts.truncateSliceOrArray > 0 {
		l, limitIfMoreThan := truncateTo(opts.truncateSliceOrArray, opts.truncateSliceOrArrayIfMoreThan)
		if n > limitIfMoreThan && n > l {
			limit = l
		}
	}

	dec = jsontext.NewDecoder(bytes.NewReader(val))
	if _, err := dec.ReadToken(); err != nil {
		return err
	}
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}
	for range limit {
		elem, err := dec.ReadValue()
		if err != nil {
			return err
		}
		if err := writeTruncatedValue(enc, elem, opts); err != nil {
			return err
		}
	}
	if limit < n && !opts.omitStatistics {
		if err := enc.WriteToken(jsontext.String(fmt.Sprintf("...%d elems", n))); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndArray)
}

// writeTruncatedObject writes the JSON object val with all its members,
// truncating the strings and arrays in their values.
func writeTruncatedObject(enc *jsontext.Encoder, val jsontext.Value, opts encOpts) error {
	dec := jsontext.NewDecoder(bytes.NewReader(val))
	if _, err := dec.ReadToken(); err != nil {
		return err
	}
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	for dec.PeekKind() != '}' {
		name, err := dec.ReadToken()
		if err != nil {
			return err
		}
		if err := enc.WriteToken(name); err != nil {
			return err
		}
		elem, err := dec.ReadValue()
		if err != nil {
			return err
		}
		if err := writeTruncatedValue(enc, elem, opts); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

// marshalTruncatedSliceOrArray writes at most limit elements from rv,
// followed by a statistics string indicating the original length.
func marshalTruncatedSliceOrArray(enc *jsontext.Encoder, rv reflect.Value, n, limit int, opts encOpts) error {
	if err := enc.WriteToken(jsontext.BeginArray); err != nil {
		return err
	}
	for i := range limit {
		if err := jsonv2.MarshalEncode(enc, rv.Index(i).Interface()); err != nil {
			return err
		}
	}
	if !opts.omitStatistics {
		if err := enc.WriteToken(jsontext.String(fmt.Sprintf("...%d elems", n))); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndArray)
}

// marshalTruncatedMap writes at most limit sorted entries from rv,
// appending a statistics suffix to the key of the next entry.
func marshalTruncatedMap(enc *jsontext.Encoder, rv reflect.Value, n, limit int, opts encOpts) error {
	type mapEntry struct {
		key    reflect.Value
		value  reflect.Value
		keyStr string
	}
	entries := make([]mapEntry, 0, n)
	for k, v := range rv.Seq2() {
		ks, err := resolveKeyName(k)
		if err != nil {
			// Let json/v2 encode (or reject) the whole map without truncation.
			return errors.ErrUnsupported
		}
		entries = append(entries, mapEntry{key: k, value: v, keyStr: ks})
	}
	slices.SortFunc(entries, func(a, b mapEntry) int {
		return strings.Compare(a.keyStr, b.keyStr)
	})

	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}

	writeLimit := limit
	var statsMsg string
	if opts.omitStatistics {
		// Just truncate, no statistics suffix.
	} else if n == limit+1 {
		// Only 1 extra — keep all to avoid wasting space on a statistics marker.
		writeLimit = n
	} else {
		// n > limit+1: truncate and append statistics to the next key.
		statsMsg = fmt.Sprintf("...%d pairs", n)
	}

	for i := 0; i < writeLimit; i++ {
		if err := enc.WriteToken(jsontext.String(entries[i].keyStr)); err != nil {
			return err
		}
		if err := jsonv2.MarshalEncode(enc, entries[i].value.Interface()); err != nil {
			return err
		}
	}

	// Write one more entry with the statistics suffix appended to its key.
	if statsMsg != "" && writeLimit < n {
		if err := enc.WriteToken(jsontext.String(entries[writeLimit].keyStr + statsMsg)); err != nil {
			return err
		}
		if err := jsonv2.MarshalEncode(enc, entries[writeLimit].value.Interface()); err != nil {
			return err
		}
	}

	return enc.WriteToken(jsontext.EndObject)
}

func stringMarshalers(opts encOpts) *jsonv2.Marshalers {
	if opts.truncateString <= 0 {
		return nil
	}

	return jsonv2.MarshalToFunc(func(enc *jsontext.Encoder, s string) error {
		// String marshalers also apply to string map keys. Never truncate an
		// object member name because doing so can change the meaning of a map.
		if isObjectName(enc) {
			return enc.WriteToken(jsontext.String(s))
		}
		return enc.WriteToken(jsontext.String(truncateStringValue(s, opts)))
	})
}

// truncateStringValue truncates s according to the string and url options.
func truncateStringValue(s string, opts encOpts) string {
	if opts.truncateString <= 0 {
		return s
	}
	var isUrl bool
	sLen := utf8.RuneCountInString(s)
	if limit, limitIfMoreThan := truncateTo(opts.truncateString, opts.truncateStringIfMoreThan); limitIfMoreThan > 0 && sLen > limitIfMoreThan {
		st := s
		if opts.truncateUrl || opts.forceLongUrl {
			u, err := url.Parse(s)
			// Require a host so that "Error: ..." or "data:..." stay plain strings.
			isUrl = err == nil && u.Scheme != "" && u.Host != ""
			if isUrl {
				if opts.forceLongUrl {
					st = s // do not truncate url, force long url
				} else {
					st = truncateUrl(s, u, opts)
				}
			}
		}
		if !isUrl {
			st = strings_.Truncate(s, limit)
			if !opts.omitStatistics {
				st += fmt.Sprintf("...%d chars", sLen)
			}
		}
		if len(st) < len(s) {
			s = st
		}
	}
	return s
}

// marshalTruncatedBytes writes b truncated as a base64 string, or returns
// [errors.ErrUnsupported] if truncating b would not make it shorter.
func marshalTruncatedBytes(enc *jsontext.Encoder, b []byte, opts encOpts) error {
	if limit, limitIfMoreThan := truncateTo(opts.truncateBytes, opts.truncateBytesIfMoreThan); limitIfMoreThan > 0 && len(b) > limitIfMoreThan {
		var m string
		if !opts.omitStatistics {
			m = fmt.Sprintf("...%d bytes", len(b))
		}
		bt := bytes_.Truncate(b, limit)
		encodedLenS := base64.StdEncoding.EncodedLen(len(b))
		encodedLenT := base64.StdEncoding.EncodedLen(len(bt))
		if len(m)+encodedLenT < encodedLenS {
			return enc.WriteToken(jsontext.String(base64.StdEncoding.EncodeToString(bt) + m))
		}
	}
	return errors.ErrUnsupported
}

// isObjectName reports whether a type-specific marshaler is currently being
// invoked for a JSON object member name. StackIndex's count is even while an
// object is expecting a name and odd while it is expecting a value.
func isObjectName(enc *jsontext.Encoder) bool {
	depth := enc.StackDepth()
	if depth <= 0 {
		return false
	}

	kind, count := enc.StackIndex(depth)
	return kind == '{' && count%2 == 0
}
