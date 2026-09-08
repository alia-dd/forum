package handlers

import (
	"errors"
	"net/http"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
)

// Likely nonsense, maybe an example of how to do it, maybe an example of how not to do it
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	w.Write([]byte(code))
	w.Write([]byte(message))
}

func handleError(w http.ResponseWriter, err error, entity string, relationships int) {
	var (
		code    string
		message string
		status  int
	)

	switch {
	case errors.Is(err, customerrors.ErrIDNotANumber):
		code = "ID_NOT_A_NUMBER"
		message = "ID is not a number"
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrNotFound):
		code = "NOT_FOUND"
		message = customerrors.ErrNotFound.Error()
		status = http.StatusNotFound
	case errors.Is(err, customerrors.ErrDuplicateEntry):
		code = "DUPLICATE_ENTRY"
		message = customerrors.ErrDuplicateEntry.Error()
		status = http.StatusConflict
	case errors.Is(err, customerrors.ErrInvalidData):
		code = "INVALID_DATA"
		message = customerrors.ErrInvalidData.Error()
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrInvalidName):
		code = "INVALID_NAME"
		message = customerrors.ErrInvalidName.Error()
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrDuplicateEmail):
		code = "DUPLICATE_EMAIL"
		message = customerrors.ErrDuplicateEmail.Error()
		status = http.StatusConflict
	case errors.Is(err, customerrors.ErrInternalError):
		code = "INTERNAL_ERROR"
		message = customerrors.ErrInternalError.Error()
		status = http.StatusInternalServerError

	default:
		code = "INTERNAL_ERROR"
		message = "An unexpected error occurred while processing the request"
		status = http.StatusInternalServerError
	}
	writeError(w, status, code, message)
}
