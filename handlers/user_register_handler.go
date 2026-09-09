package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/service"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

type UseHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UseHandler {
	return &UseHandler{service: service}
}

func (h *UseHandler) GetRegisterUser(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, http.StatusOK, "user_registration", nil)
}

func (h *UseHandler) PostRegisterUser(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()

	if parseErr := r.ParseForm(); parseErr != nil {
		handleError(w, customerrors.ErrInternalError)
		return
	}

	userData := models.UserRegister{
		Username:        strings.TrimSpace(r.FormValue("username")),
		Name:            strings.TrimSpace(r.FormValue("fullname")),
		Email:           strings.TrimSpace(r.FormValue("email")),
		Password:        strings.TrimSpace(r.FormValue("password")),
		ConformPassword: strings.TrimSpace(r.FormValue("confirm_password")),
	}

	if PostErr := h.service.CreateUserService(cx, userData); PostErr != nil {
		handleError(w, PostErr)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// this does a server side extra check on the usename and email for duplicets and incorrect format
func (h *UseHandler) CheckIfAvailabe(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	var (
		exists   bool
		checkErr error
		message  string
	)
	switch r.URL.Path {
	case "/check-username":
		username := r.URL.Query().Get("username")
		if !utils.IsValidName(username) {
			message = customerrors.ErrInvalidName.Error()
			break
		}
		exists, checkErr = h.service.CheckIfAvailable(cx, username)
		message = fmt.Sprintf("Username %s is not availabl", username)
	case "/check-email":
		email := r.URL.Query().Get("email")
		if !utils.IsValidEmail(email) {
			message = customerrors.ErrInvalidData.Error()
			break
		}
		exists, checkErr = h.service.CheckIfAvailable(cx, email)
		message = customerrors.ErrDuplicateEmail.Error()
	default:
		handleError(w, customerrors.ErrInternalError)
		return
	}
	if checkErr != nil {
		http.Error(w, "Failed to Check Availability", http.StatusInternalServerError)
		return
	}

	payload := struct {
		Exists  bool   `json:"exists"`
		Message string `json:"message"`
	}{
		Exists:  exists,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}
