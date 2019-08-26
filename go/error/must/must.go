package must

// Must panics if err != nil
func Must(err error) {
	if err == nil {
		return
	}
	panic(err)
}
