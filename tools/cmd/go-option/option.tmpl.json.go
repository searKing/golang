package main

const tmplOption = `
// A {{.OptionInterfaceName}} sets options.
type {{.OptionInterfaceName}} interface {
	apply(*{{.OptionTypeName}})
}

// Empty{{.OptionInterfaceName}} does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type Empty{{.OptionInterfaceName}} struct{}

func (Empty{{.OptionInterfaceName}}) apply(*{{.OptionTypeName}}) {}

// {{.OptionInterfaceName}}Func wraps a function that modifies {{.OptionTypeName}} into an
// implementation of the {{.OptionInterfaceName}} interface.
type {{.OptionInterfaceName}}Func func(*{{.OptionTypeName}})

func (f {{.OptionInterfaceName}}Func) apply(do *{{.OptionTypeName}}) {
	f(do)
}

// sample code for option, default for nothing to change
func _{{.OptionInterfaceName}}WithDefault() {{.OptionInterfaceName}} {
	return {{.OptionInterfaceName}}Func(func( *{{.OptionTypeName}}) {
		// nothing to change
	})
}

{{- if .ApplyOptionsAsMemberFunction }}
func (o *{{.OptionTypeName}}) ApplyOptions(options ...{{.OptionInterfaceName}}) *{{.OptionTypeName}} {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}

type completed{{.OptionTypeName}} struct {
	*{{.OptionTypeName}}
}

type Completed{{.OptionTypeName}} struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completed{{.OptionTypeName}}
}
// recommend to implement codes below
// Code borrowed from https://github.com/kubernetes/kubernetes

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// func (o *{{.OptionTypeName}}) Complete() Completed{{.OptionTypeName}} {
//  o.ApplyOptions(options...)
//  // Add custom codes here
//  return Completed{{.OptionTypeName}}{&completed{{.OptionTypeName}}{o}}
// }


// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
//func (c completed{{.OptionTypeName}}) New(name string) (*{{.OptionTypeName}}, error) {
//  // Add custom codes here
//	return nil, fmt.Errorf("not implemented")
//}

// Apply set options and something else as global init, act likes New but without {{.OptionTypeName}}'s instance
// Apply usually called after Complete
//func (c completed{{.OptionTypeName}}) Apply() error {
//  // Add custom codes here
//	return fmt.Errorf("not implemented")
//}

{{- else}}
func ApplyOptions(o *{{.OptionTypeName}}, options ...{{.OptionInterfaceName}}) *{{.OptionTypeName}} {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
{{- end}}
`
