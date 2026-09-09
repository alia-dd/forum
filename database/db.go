package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data/forum.db?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if errPing := db.Ping(); errPing != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to ping database: %w", errPing)
	}
	_, err = InitTable(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}
