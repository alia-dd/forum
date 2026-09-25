package handlers

import (
	"errors"
	"fmt"
	"net/http"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		//code    string
		message string
		status  int
	)

	switch {
	case errors.Is(err, customerrors.ErrNotFound):
		//code = "NOT_FOUND"
		message = "That item does not exist"
		status = http.StatusNotFound
	case errors.Is(err, customerrors.ErrForeignKeyConstraint):
		//code = "INVALID_REFERENCE"
		message = "Referenced item does not exist"
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrNotFound):
		//code = "NOT_FOUND"
		message = customerrors.ErrNotFound.Error()
		status = http.StatusNotFound
	case errors.Is(err, customerrors.ErrDuplicateEntry):
		//code = "DUPLICATE_ENTRY"
		message = customerrors.ErrDuplicateEntry.Error()
		status = http.StatusConflict
	case errors.Is(err, customerrors.ErrInvalidData):
		//code = "INVALID_DATA"
		message = customerrors.ErrInvalidData.Error()
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrInvalidName):
		//code = "INVALID_NAME"
		message = customerrors.ErrInvalidName.Error()
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrDuplicateEmail):
		//code = "DUPLICATE_EMAIL"
		message = customerrors.ErrDuplicateEmail.Error()
		status = http.StatusConflict
	case errors.Is(err, customerrors.ErrInternalError):
		//code = "INTERNAL_ERROR"
		message = customerrors.ErrInternalError.Error()
		status = http.StatusInternalServerError

	case errors.Is(err, customerrors.ErrBadRequest):
		//code = "BAD_REQUEST"
		message = "Invalid request"
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrForbidden):
		//code = "NOT_ALLOWED"
		message = "You cannot do that"
		status = http.StatusForbidden
	case errors.Is(err, customerrors.ErrDuplicateEntry):
		//code = "DUPLICATE_ENTRY"
		message = "That item already exists"
		status = http.StatusConflict
	case errors.Is(err, customerrors.ErrInUse):
		//code = "IN_USE"
		message = "This item is still referenced and cannot be deleted"
		status = http.StatusConflict
	case errors.Is(err, customerrors.ErrIncorrectPassword):
		// code = "INCORRECT_PASSWORD"
		message = customerrors.ErrIncorrectPassword.Error()
		status = http.StatusBadRequest
	case errors.Is(err, customerrors.ErrInvalidLogin):
		// code = "INVALID_LOGIN"
		message = customerrors.ErrInvalidLogin.Error()
		status = http.StatusBadRequest
	default:
		//code = "INTERNAL_ERROR"
		message = "An unexpected error occurred while processing the request"
		status = http.StatusInternalServerError
	}

	user, _ := r.Context().Value("user_session").(*models.UserInfo) //nil is fine, just to show top banner correctly

	pageData := models.PageData[models.ErrorStruct]{
		User: user,
		PageContent: models.ErrorStruct{
			Error:   fmt.Sprint(status),
			ErrorMs: message,
		},
	}

	utils.RenderTemplate(w, status, "error", pageData)
}
