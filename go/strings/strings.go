package strings

// SliceContains  reports whether s is within ss.
func SliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
