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
	ErrIncorrectPassword    = errors.New("Incorrect Password")
	ErrInvalidLogin         = errors.New("Invalid username or password")
	ErrBlockedUser          = errors.New("Your account has been blocked")
	ErrInvalidImageType     = errors.New("Image must be a JPEG or PNG.")
	ErrImageTooLarge        = errors.New("Image must be less than 10 MB.")
	ErrCommentEmpty         = errors.New("Comment cannot be empty.")
	ErrCommentTooLong       = errors.New("Comment has to be less than 2,000 characters long.")
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
