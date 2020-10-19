package strings

// Return value if nonempty, def otherwise.
func ValueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}
