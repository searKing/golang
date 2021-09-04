// Code generated by "go-option -type resolveAddr"; DO NOT EDIT.

package resolver

// A ResolveAddrOption sets options.
type ResolveAddrOption interface {
	apply(*resolveAddr)
}

// EmptyResolveAddrOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyResolveAddrOption struct{}

func (EmptyResolveAddrOption) apply(*resolveAddr) {}

// ResolveAddrOptionFunc wraps a function that modifies resolveAddr into an
// implementation of the ResolveAddrOption interface.
type ResolveAddrOptionFunc func(*resolveAddr)

func (f ResolveAddrOptionFunc) apply(do *resolveAddr) {
	f(do)
}

// sample code for option, default for nothing to change
func _ResolveAddrOptionWithDefault() ResolveAddrOption {
	return ResolveAddrOptionFunc(func(*resolveAddr) {
		// nothing to change
	})
}
func (o *resolveAddr) ApplyOptions(options ...ResolveAddrOption) *resolveAddr {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}