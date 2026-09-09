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
	ErrDatabaseBusy         = errors.New("Database in use by another user or process")
	ErrIDNotANumber         = errors.New("Id not a number")
	ErrInvalidData          = errors.New("Invalid Entery")
	ErrInternalError        = errors.New("Internal Server Error")
	ErrInvalidName          = errors.New("Username may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen.")
	ErrDuplicateEmail       = errors.New("The email you have provided is already associated with an account.")
	ErrBadRequest           = errors.New("bad request")
	ErrForbidden            = errors.New("forbidden")
	ErrInUse                = errors.New("resource still in use")
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
