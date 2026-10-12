package database

import (
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var schema string

func InitSchema(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
