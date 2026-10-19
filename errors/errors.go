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
	ErrInvalidData          = errors.New("Invalid Entry")
	ErrInternalError        = errors.New("Internal Server Error")
	ErrInvalidName          = errors.New("Username must be between 3 and 20 characters long and may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen.")
	ErrDuplicateEmail       = errors.New("The email you have provided is already associated with an account.")
	ErrInvalidEmail         = errors.New("The email you have provided is invalid.")
	ErrBadRequest           = errors.New("Bad request")
	ErrForbidden            = errors.New("Forbidden")
	ErrInUse                = errors.New("Resource still in use")
	ErrInvalidBio           = errors.New("Bio must be 500 characters or less.")
	ErrInvalidImageType     = errors.New("Image must be a JPEG or PNG.")
	ErrImageTooLarge        = errors.New("Image must be less than 10 MB.")
	ErrCommentEmpty         = errors.New("Comment cannot be empty.")
	ErrCommentTooLong       = errors.New("Comment has to be less than 3,000 characters long.")
	ErrIncorrectPassword    = errors.New("Incorrect password.")
	ErrInvalidLogin         = errors.New("Invalid username or password.")
	ErrBlockedUser          = errors.New("Your account has been blocked.")
	ErrDuplicateUsername    = errors.New("The username you have provided is already associated with an account.")
	ErrPasswordMismatch     = errors.New("Password does not match confirmation password.")
	ErrUsernameWrongLength  = errors.New("Invalid username: it must be 3 to 50 characters long.")
	ErrPasswordLength       = errors.New("Password myst be between 8 and 100 characters.")
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
