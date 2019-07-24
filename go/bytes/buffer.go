package bytes

import "bytes"

func NewIndent(dst *bytes.Buffer, prefix, indent string, depth int) {
	dst.WriteString(prefix)
	for i := 0; i < depth; i++ {
		dst.WriteString(indent)
	}
}

func NewLine(dst *bytes.Buffer, prefix, indent string, depth int) {
	dst.WriteByte('\n')
	NewIndent(dst, prefix, indent, depth)
}
