package strings

// ValueOrDefault Return first value nonempty
// Example:
//	ValueOrDefault(value, def)
func ValueOrDefault(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
