package context

type nopTags struct{}

func (t *nopTags) Set(key string, value interface{}) {
	return
}

func (t *nopTags) Get(key string) (interface{}, bool) {
	return nil, false
}

// Del deletes the values associated with key.
func (t *nopTags) Del(key string) {
	return
}
func (t *nopTags) Values() map[string]interface{} {
	return nil
}
