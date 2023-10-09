package api

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

type Database struct {
	ConnectionString string
	Database         *sql.DB
}

func LoginDatabase(login, password string) (db *Database, err error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	connStr := "postgres://" + login + ":" + password + "@" + host + ":" + port + "/" + name + "?sslmode=disable"

	db, e := NewDatabase(connStr)
	if e != nil {
		return nil, e
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDatabase(connectionString string) (db *Database, err error) {
	database, e := sql.Open("postgres", connectionString)
	if e != nil {
		return nil, e
	}
	return &Database{
		ConnectionString: connectionString,
		Database:         database,
	}, nil
}

func (db *Database) Disconnect() (err error) {
	return db.Database.Close()
}

func (db *Database) Ping() (err error) {
	return db.Database.Ping()
}
