package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (h *UseHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	statusCode := 200

	if r.URL.RawQuery != "" {
		var userData models.UserRegister
		jsonErr := json.NewDecoder(r.Body).Decode(&userData)
		if jsonErr != nil {
			http.Error(w, "incorrect data", http.StatusBadRequest)
			return
		}
		if PostErr := h.service.CreateUserService(cx, userData); PostErr != nil {
			http.Error(w, "something here dont know what mmmmmm :/", http.StatusInternalServerError)
		}
		statusCode = 201
	}

	utils.RenderTemplate(w, statusCode, "user_registration", nil)
}
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
		exists, checkErr = h.service.CheckIfAvailable(cx, username)
		message = fmt.Sprintf("Username %s is not availabl", username)
	case "/check-email":
		email := r.URL.Query().Get("email")
		exists, checkErr = h.service.CheckIfAvailable(cx, email)
		message = "The email you have provided is already associated with an account."
	default:
		http.Error(w, "not found", http.StatusNotFound)
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

func (h *UseHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	userId, idErr := strconv.Atoi(r.PathValue("id"))
	if idErr != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if GetErr := h.service.GetUserByIdService(cx, userId); GetErr != nil {

	}
}
