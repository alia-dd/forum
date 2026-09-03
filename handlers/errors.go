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
	default:
		code = "INTERNAL_ERROR"
		message = "An unexpected error occurred while processing the request"
		status = http.StatusInternalServerError
	}

	writeError(w, status, code, message)
}
