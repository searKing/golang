// Code generated by "go-option -type=option --trim"; DO NOT EDIT.
// Install go-option by "go get install github.com/searKing/golang/tools/go-option"

package otlptracegrpc

import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

// A Option sets options.
type Option interface {
	apply(*option)
}

// EmptyOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyOption struct{}

func (EmptyOption) apply(*option) {}

// OptionFunc wraps a function that modifies option into an
// implementation of the Option interface.
type OptionFunc func(*option)

func (f OptionFunc) apply(do *option) {
	f(do)
}

// ApplyOptions call apply() for all options one by one
func (o *option) ApplyOptions(options ...Option) *option {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}

// withOption sets option.
func withOption(v option) Option {
	return OptionFunc(func(o *option) {
		*o = v
	})
}

// WithOptionOtlpOptions appends OtlpOptions in option.
func WithOptionOtlpOptions(v ...otlptracegrpc.Option) Option {
	return OptionFunc(func(o *option) {
		o.OtlpOptions = append(o.OtlpOptions, v...)
	})
}

// WithOptionOtlpOptionsReplace sets OtlpOptions in option.
func WithOptionOtlpOptionsReplace(v ...otlptracegrpc.Option) Option {
	return OptionFunc(func(o *option) {
		o.OtlpOptions = v
	})
}