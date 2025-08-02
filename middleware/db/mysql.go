package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/zhouximi/email_compromise_checker/data_model"
	"github.com/zhouximi/email_compromise_checker/types"
	"os"
)

var filePath = "config/mysql_config.json"

type MySQLQuerier struct {
	db *sql.DB
}

func NewMySQLQuerier() (*MySQLQuerier, error) {
	dbConfig, err := loadDBConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)

	// Open DB connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Set connection pool options (optional but recommended)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	// Ping DB to verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &MySQLQuerier{db: db}, nil
}

func (q *MySQLQuerier) RunQuery(key string, queryFn func(db *sql.DB, key string) (interface{}, error)) (interface{}, error) {
	return queryFn(q.db, key)
}

func loadDBConfig() (*data_model.MySQLConfig, error) {
	dbConfig, err := os.ReadFile(filePath)
	if err != nil {
		return nil, types.ErrReadConfigFile
	}

	var cfg *data_model.MySQLConfig
	if err := json.Unmarshal(dbConfig, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
