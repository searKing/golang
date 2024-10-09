// Code generated by "go-option -type=option --trim"; DO NOT EDIT.
// Install go-option by "go get install github.com/searKing/golang/tools/go-option"

package prometheusmetric

import (
	prometheusmetric "go.opentelemetry.io/otel/exporters/prometheus"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

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

// WithOptionPrometheusOptions appends PrometheusOptions in option.
func WithOptionPrometheusOptions(v ...prometheusmetric.Option) Option {
	return OptionFunc(func(o *option) {
		o.PrometheusOptions = append(o.PrometheusOptions, v...)
	})
}

// WithOptionPrometheusOptionsReplace sets PrometheusOptions in option.
func WithOptionPrometheusOptionsReplace(v ...prometheusmetric.Option) Option {
	return OptionFunc(func(o *option) {
		o.PrometheusOptions = v
	})
}

// WithOptionExporterWrappers appends ExporterWrappers in option.
func WithOptionExporterWrappers(v ...func(exporter sdkmetric.Exporter) sdkmetric.Exporter) Option {
	return OptionFunc(func(o *option) {
		o.ExporterWrappers = append(o.ExporterWrappers, v...)
	})
}

// WithOptionExporterWrappersReplace sets ExporterWrappers in option.
func WithOptionExporterWrappersReplace(v ...func(exporter sdkmetric.Exporter) sdkmetric.Exporter) Option {
	return OptionFunc(func(o *option) {
		o.ExporterWrappers = v
	})
}
