package context

type mapTags struct {
	values map[string]interface{}
}

func (t *mapTags) Set(key string, value interface{}) Tags {
	t.values[key] = value
	return t
}

func (t *mapTags) Has(key string) bool {
	_, ok := t.values[key]
	return ok
}

func (t *mapTags) Values() map[string]interface{} {
	return t.values
}
