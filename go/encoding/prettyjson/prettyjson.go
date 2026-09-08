package prettyjson

import (
	"encoding/base64"
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"errors"
	"fmt"
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
	var opts encOpts
	opts.ApplyOptions(options...)

	marshalers := AdvancedTruncateMarshalers(opts)

	var jsonOpts []jsonv2.Options
	if marshalers != nil {
		jsonOpts = append(jsonOpts, jsonv2.WithMarshalers(marshalers))
	}
	jsonOpts = append(jsonOpts, jsontext.EscapeForHTML(opts.escapeHTML))
	jsonOpts = append(jsonOpts, jsonv2.Deterministic(true))

	if opts.omitEmpty {
		// encoding/json's v1 compatibility option is an encoding/json/v2
		// Options value and can therefore be passed directly to json/v2.
		jsonOpts = append(jsonOpts, jsonv1.OmitEmptyWithLegacySemantics(true))
		// OmitZeroStructFields extends omit-empty to zero-value structs,
		// matching the custom isEmptyValue/v.IsZero() behavior of the old encoder.
		jsonOpts = append(jsonOpts, jsonv2.OmitZeroStructFields(true))
	}

	return jsonv2.Marshal(v, jsonOpts...)
}

// AdvancedTruncateMarshalers returns the type-specific marshalers used by
// Marshal. Concrete types (string, []byte) are handled by dedicated marshalers.
// Arbitrary slice/array/map values are handled by the catch-all "any" marshaler
// which uses reflect.Kind to distinguish maps from structs — something the
// previous jsontext post-processing pass (truncateContainers) could not do.
func AdvancedTruncateMarshalers(opts encOpts) *jsonv2.Marshalers {
	return jsonv2.JoinMarshalers(
		stringMarshalers(opts),
		bytesMarshalers(opts),
		anyMarshalers(opts),
	)
}

func anyMarshalers(opts encOpts) *jsonv2.Marshalers {
	if opts.truncateSliceOrArray <= 0 && opts.truncateMap <= 0 {
		return nil
	}

	return jsonv2.MarshalToFunc(func(enc *jsontext.Encoder, v any) error {
		rv := reflect.ValueOf(v)
		if !rv.IsValid() {
			return errors.ErrUnsupported
		}
		// json/v2 passes an addressable pointer to the value; dereference it.
		for rv.Kind() == reflect.Pointer {
			if rv.IsNil() {
				return errors.ErrUnsupported
			}
			rv = rv.Elem()
		}

		switch rv.Kind() {
		case reflect.Slice:
			if opts.truncateSliceOrArray <= 0 || rv.IsNil() {
				return errors.ErrUnsupported
			}
			// Skip byte slices — the default marshaler base64-encodes them.
			if rv.Type().Elem().Kind() == reflect.Uint8 {
				return errors.ErrUnsupported
			}
			n := rv.Len()
			limit, limitIfMoreThan := truncateTo(opts.truncateSliceOrArray, opts.truncateSliceOrArrayIfMoreThan)
			if n <= limitIfMoreThan || n <= limit {
				return errors.ErrUnsupported
			}
			return marshalTruncatedSliceOrArray(enc, rv, n, limit, opts)

		case reflect.Array:
			if opts.truncateSliceOrArray <= 0 {
				return errors.ErrUnsupported
			}
			n := rv.Len()
			limit, limitIfMoreThan := truncateTo(opts.truncateSliceOrArray, opts.truncateSliceOrArrayIfMoreThan)
			if n <= limitIfMoreThan || n <= limit {
				return errors.ErrUnsupported
			}
			return marshalTruncatedSliceOrArray(enc, rv, n, limit, opts)

		case reflect.Map:
			if opts.truncateMap <= 0 || rv.IsNil() {
				return errors.ErrUnsupported
			}
			n := rv.Len()
			limit, limitIfMoreThan := truncateTo(opts.truncateMap, opts.truncateMapIfMoreThan)
			if n <= limitIfMoreThan || n <= limit {
				return errors.ErrUnsupported
			}
			return marshalTruncatedMap(enc, rv, n, limit, opts)

		default:
			// Structs and all other types use the default marshaler.
			return errors.ErrUnsupported
		}
	})
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
			return fmt.Errorf("json: unsupported type: %s", rv.Type().Key())
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
		var isUrl bool
		sLen := utf8.RuneCountInString(s)
		if limit, limitIfMoreThan := truncateTo(opts.truncateString, opts.truncateStringIfMoreThan); limitIfMoreThan > 0 && sLen > limitIfMoreThan {
			st := s
			if opts.truncateUrl || opts.forceLongUrl {
				u, err := url.Parse(s)
				isUrl = err == nil && u.Scheme != ""
				if isUrl {
					if opts.forceLongUrl {
						st = s // do not truncate url, force long url
					} else {
						st = truncateUrl(u, opts)
					}
				}
			}
			if !isUrl {
				st = strings_.Truncate(s, limit)
				if !opts.omitStatistics {
					st += fmt.Sprintf("...%d chars", len(s))
				}
			}
			if len(st) < len(s) {
				s = st
			}
		}
		return enc.WriteToken(jsontext.String(s))
	})
}

func bytesMarshalers(opts encOpts) *jsonv2.Marshalers {
	if opts.truncateBytes <= 0 {
		return nil
	}

	return jsonv2.MarshalFunc(func(b []byte) ([]byte, error) {
		if limit, limitIfMoreThan := truncateTo(opts.truncateBytes, opts.truncateBytesIfMoreThan); limitIfMoreThan > 0 && len(b) > limitIfMoreThan {
			var m string
			if !opts.omitStatistics {
				m = fmt.Sprintf("...%d bytes", len(b))
			}
			bt := bytes_.Truncate(b, limit)
			encodedLenS := base64.StdEncoding.EncodedLen(len(b))
			encodedLenT := base64.StdEncoding.EncodedLen(len(bt))
			if len(m)+encodedLenT < encodedLenS {
				return jsontext.AppendQuote(nil, base64.StdEncoding.EncodeToString(bt)+m)
			}
		}

		return jsonv2.Marshal(b)
	})
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
