package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLQuerier struct {
	db *sql.DB
}

func NewMySQLQuerier() (*MySQLQuerier, error) {
	db, err := sql.Open("mysql", getDNSAddress())
	if err != nil {
		return nil, err
	}

	// Set connection pool options (optional but recommended)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	// Ping DB to verify connection
	if err := db.Ping(); err != nil {
		log.Fatal("[NewMySQLQuerier]failed to connect db")
		return nil, err
	}

	return &MySQLQuerier{db: db}, nil
}

func (q *MySQLQuerier) RunQuery(key string, queryFn func(db *sql.DB, key string) (interface{}, error)) (interface{}, error) {
	return queryFn(q.db, key)
}

func getDNSAddress() string {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_DATABASE")

	if user == "" || pass == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("Missing one or more required environment variables")
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
}
