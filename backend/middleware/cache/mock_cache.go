package cache

type MockCache struct {
	Items map[string]interface{}
}

func (m *MockCache) Get(cacheKey string) (interface{}, error) {
	if value, exists := m.Items[cacheKey]; exists {
		return value, nil
	}
	return nil, nil
}

func (m *MockCache) Set(cacheKey string, value interface{}) error {
	if m.Items == nil {
		m.Items = make(map[string]interface{})
	}
	m.Items[cacheKey] = value
	return nil
}
