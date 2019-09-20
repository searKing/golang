package context

type nopTags struct{}

func (t *nopTags) Set(key string, value interface{}) Tags {
	return t
}

func (t *nopTags) Get(key string) (interface{}, bool) {
	return nil, false
}

func (t *nopTags) Values() map[string]interface{} {
	return nil
}
