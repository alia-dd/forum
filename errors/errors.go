package customerrors

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

// Example of functionality, probably not using what exists currently,
// just there to show how it was done last time

var (
	//messages here are for internal use
	// messages to be sent externally are configured in handlers/errors.go
	ErrNotFound             = errors.New("resource not found")
	ErrDuplicateEntry       = errors.New("entry already exists")
	ErrForeignKeyConstraint = errors.New("Foreign key error")
	ErrBadRequest           = errors.New("bad request")
	ErrForbidden            = errors.New("forbidden")
)

func MapSQLError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		switch sqliteErr.ExtendedCode {
		case sqlite3.ErrConstraintUnique:
			return ErrDuplicateEntry
		case sqlite3.ErrConstraintForeignKey:
			return ErrForeignKeyConstraint
		}

	}

	return err
}
