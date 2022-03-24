package main

type TmplPackageRender struct {
	// the header
	GenerateToolName string
	GenerateToolArgs string

	// package clause
	PackageName string // Package name

	// sync.Map clause.
	MapRenders []TmplMapRender

	// Options
	WithSyncMapMethod       bool
	WithMethodLoadAndDelete bool
}

type TmplMapRender struct {
	MapTypeName string // Map Type Name of the sync.Map type.
	MapImport   string // Import path of the sync.Map type.

	KeyTypeName   string // Key Type of the sync.Map type.
	KeyTypeImport string // Import path of the sync.Map's key.

	ValueTypeName    string // Value Type of the sync.Map type.
	ValueTypeImport  string // Import path of the sync.Map's Value.
	ValueTypeNilVal  string // Nil value of the sync.Map's Value.
	ValueTypeNilDecl string // Nil decl of the sync.Map's Value.
}
