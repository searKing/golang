package reflect

import (
	"fmt"
	"strings"
)

// SubStructTag defines a single struct's string literal tag
type SubStructTag struct {
	// Key is the tag key, such as json, xml, etc.
	// i.e: `json:"foo,omitempty". Here key is: "json"
	Key string

	// Name is a part of the value
	// i.e: `json:"foo,omitempty". Here name is: "foo"
	Name string

	// Options is a part of the value. It contains a slice of tag options i.e:
	// `json:"foo,omitempty". Here options is: ["omitempty"]
	Options []string
}

// HasOption returns true if the given option is available in options
func (t *SubStructTag) HasOption(opt string) bool {
	return hasOption(t.Options, opt)
}

// AddOptions adds the given option.
// It ignores any duplicated option.
func (t *SubStructTag) AddOptions(opts ...string) {
	for _, opt := range opts {
		if t.HasOption(opt) {
			continue
		}
		t.Options = append(t.Options, opt)
	}

	return
}

// DeleteOptions removes the given option.
// It ignores any option not found.
func (t *SubStructTag) DeleteOptions(opts ...string) {
	var cleanOpts []string
	for _, opt := range t.Options {
		if hasOption(opts, opt) {
			continue
		}
		cleanOpts = append(cleanOpts, opt)
	}

	t.Options = cleanOpts
	return
}

// Value returns the raw value of the tag, i.e. if the tag is
// `json:"foo,omitempty", the Value is "foo,omitempty"
func (t *SubStructTag) Value() string {
	options := strings.Join(t.Options, ",")
	if options != "" {
		return fmt.Sprintf(`%s,%s`, t.Name, options)
	}
	return t.Name
}

// String reassembles the tag into a valid tag field representation
func (t *SubStructTag) String() string {
	return fmt.Sprintf(`%s:%q`, t.Key, t.Value())
}

// GoString implements the fmt.GoStringer interface
func (t *SubStructTag) GoString() string {
	return fmt.Sprintf(`{
		Key:    '%s',
		Name:   '%s',
		Option: '%s',
	}`, t.Key, t.Name, strings.Join(t.Options, ","))
}

func hasOption(options []string, opt string) bool {
	for _, tagOpt := range options {
		if tagOpt == opt {
			return true
		}
	}
	return false
}
