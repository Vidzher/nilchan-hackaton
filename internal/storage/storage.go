package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			balance INTEGER NOT NULL DEFAULT 0
		);
		CREATE INDEX IF NOT EXISTS idx_email on users(email)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %s", err.Error())
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %s", err.Error())
	}

	return &Storage{DB: db}, nil
}
