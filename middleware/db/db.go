package db

import "database/sql"

type IDB interface {
	RunQuery(key string, queryFn func(db *sql.DB, key string) (interface{}, error)) (interface{}, error)
}
