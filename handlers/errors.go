package handlers

import (
	"errors"
	"net/http"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
)

// needs to be changed to write/utilize error.html template or somesuch
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	w.Write([]byte(code))
	w.Write([]byte(message))
}

func handleError(w http.ResponseWriter, err error) {
	var (
		code    string
		message string
		status  int
	)

	switch {
	case errors.Is(err, customerrors.ErrNotFound):
		code = "NOT_FOUND"
		message = "That item does not exist"
		status = http.StatusNotFound
	case errors.Is(err, customerrors.ErrForeignKeyConstraint):
		code = "INVALID_REFERENCE"
		message = "Referenced item does not exist"
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrBadRequest):
		code = "BAD_REQUEST"
		message = "Invalid request"
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrForbidden):
		code = "NOT_ALLOWED"
		message = "You cannot do that"
		status = http.StatusForbidden
	case errors.Is(err, customerrors.ErrDuplicateEntry):
		code = "DUPLICATE_ENTRY"
		message = "That item already exists"
		status = http.StatusConflict
	default:
		code = "INTERNAL_ERROR"
		message = "An unexpected error occurred while processing the request"
		status = http.StatusInternalServerError
	}

	writeError(w, status, code, message)
}
