package customerrors

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound             = errors.New("resource not found")
	ErrDuplicateEntry       = errors.New("entry already exists")
	ErrForeignKeyConstraint = errors.New("Foreign key error")
	ErrDatabaseBusy         = errors.New("Database in use by another user or process")
	ErrIDNotANumber         = errors.New("ID is not a number")
	ErrInvalidData          = errors.New("Invalid Entery")
	ErrInternalError        = errors.New("Internal Server Error")
	ErrInvalidName          = errors.New("Username may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen.")
	ErrDuplicateEmail       = errors.New("The email you have provided is already associated with an account.")
	ErrInvalidEmail         = errors.New("The email you have provided is invalid.")
	ErrBadRequest           = errors.New("Bad request")
	ErrForbidden            = errors.New("Forbidden")
	ErrInUse                = errors.New("Resource still in use")
	ErrIncorrectPassword    = errors.New("Invalid password: it must be 8 to 100 characters.")
	ErrInvalidLogin         = errors.New("Invalid username or password.")
	ErrBlockedUser          = errors.New("Your account has been blocked.")
	ErrDuplicateUsername    = errors.New("The username you have provided is already associated with an account.")
	ErrPasswordMismatch     = errors.New("password does not match confirmation password.")
	ErrUsernameWrongLength  = errors.New("Username must be between 3 and 20 characters long and.")
	ErrInvalidBio           = errors.New("Bio must be 500 characters or less.")
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
