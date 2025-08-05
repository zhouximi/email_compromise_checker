package db

import (
	"database/sql"
	"github.com/zhouximi/email_compromise_checker/data_model"
)

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
	m.Items[cacheKey] = value
	return nil
}

type MockDB struct {
	CompromisedEmails map[string]bool
}

func (m *MockDB) RunQuery(key string, queryFn func(*sql.DB, string) (interface{}, error)) (interface{}, error) {
	if compromised, exists := m.CompromisedEmails[key]; exists {
		return &data_model.EmailInfo{Email: key, Compromised: compromised}, nil
	}
	return &data_model.EmailInfo{Email: key, Compromised: false}, nil
}
