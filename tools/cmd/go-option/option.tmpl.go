package main

import "strings"

type TmplRender struct {
	OptionInterfaceName string // Option Interface Name for Option Type
	OptionTypeName      string // Option Type Name
	OptionTypeImport    string // Option Import Path

	ApplyOptionsAsMemberFunction bool // ApplyOptions can be registered as OptionType's member function
}

func (t *TmplRender) Complete() {
	t.ApplyOptionsAsMemberFunction = strings.TrimSpace(t.OptionTypeImport) == ""
}
