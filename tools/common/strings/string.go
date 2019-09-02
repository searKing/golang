package strings

func LoadElse(loaded bool, loadValue string, elseValue string) string {
	if loaded {
		return loadValue
	}
	return elseValue
}

func LoadElseGet(loaded bool, loadValue string, elseValueGetter func() string) string {
	if loaded {
		return loadValue
	}
	if elseValueGetter == nil {
		return ""
	}
	return elseValueGetter()
}
