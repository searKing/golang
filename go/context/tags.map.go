package context

type mapTags struct {
	values map[string]interface{}
}

func (t *mapTags) Set(key string, value interface{}) Tags {
	t.values[key] = value
	return t
}

func (t *mapTags) Get(key string) (interface{}, bool) {
	val, ok := t.values[key]
	return val, ok
}

func (t *mapTags) Values() map[string]interface{} {
	return t.values
}
