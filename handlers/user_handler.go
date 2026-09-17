package handlers

import (
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

	r.ParseMultipartForm(10 << 20)
	// change the handleErr to return to htmx page instead of new page
	if parseErr := r.ParseForm(); parseErr != nil {
		handleError(w, customerrors.ErrInternalError)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	name := strings.TrimSpace(r.FormValue("fullname"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))
	conformPassword := strings.TrimSpace(r.FormValue("confirm_password"))

	file, header, fileErr := r.FormFile("imagepath")
	if fileErr != nil || SaveUploadedFile(file, header.Filename) != nil {
		pageData := models.PageData{
			User:    nil,
			IsOwner: false,
			PageContent: &models.UserRegister{
				Username:        username,
				Name:            name,
				Email:           email,
				Password:        password,
				ConformPassword: conformPassword,
			},
			Error: customerrors.ErrBadRequest.Error(),
		}
		utils.RenderTemplate(w, http.StatusBadRequest, "user_registration", pageData)

	}
	fileName := header.Filename

	defer file.Close()

	userData := models.UserRegister{
		Username:        username,
		Name:            name,
		Email:           email,
		Image:           fileName,
		Password:        password,
		ConformPassword: conformPassword,
	}

	if PostErr := h.service.CreateUserService(cx, userData); PostErr != nil {
		pageData := models.PageData{
			User:    nil,
			IsOwner: false,
			PageContent: &models.UserRegister{
				Username:        username,
				Name:            name,
				Email:           email,
				Password:        password,
				ConformPassword: conformPassword,
			},
			Error: PostErr.Error(),
		}
		utils.RenderTemplate(w, http.StatusBadRequest, "user_registration", pageData)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *UseHandler) GetEditUserProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("is it working?")
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	pageData := models.PageData{
		User:        user,
		IsOwner:     ok,
		PageContent: nil,
	}

	utils.RenderTemplate(w, http.StatusOK, "profileEdit", pageData)
}

// shows profile data for other registred users
func (h *UseHandler) GetOtherUserProfile(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()

	_, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	username := r.PathValue("username")
	fmt.Println("name", username)
	profileData, fetchErr := h.service.GetUserService(cx, username)
	if fetchErr != nil {
		handleError(w, fetchErr)
		return
	}

	fmt.Println(profileData)
	// is this actually secure i hope so
	// needs extra fix
	// new problem if i were to go to  my profile why clicking on the auther name(me)
	// the user will not be able to edit or change their profile is that actual an edge case or what
	pageData := models.PageData{
		User:        &profileData,
		IsOwner:     false,
		PageContent: nil,
	}
	utils.RenderTemplate(w, http.StatusOK, "profile", pageData)
}

func (h *UseHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()

	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	if parseErr := r.ParseForm(); parseErr != nil {
		handleError(w, customerrors.ErrInternalError)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	name := strings.TrimSpace(r.FormValue("fullname"))
	email := strings.TrimSpace(r.FormValue("email"))

	userData := models.UserUpdate{
		Id: user.Id,
	}
	if username != "" {
		userData.Username = &username
	}
	if name != "" {
		userData.Name = &name
	}
	if email != "" {
		userData.Email = &email
	}

	if UpdateErr := h.service.UpdateUserService(cx, userData); UpdateErr != nil {
		pageData := models.PageData{
			User: &models.UserInfo{
				Id:       user.Id,
				Username: username,
				Name:     name,
				Email:    email,
			},
			IsOwner:     true,
			PageContent: nil,
			Error:       UpdateErr.Error(),
		}
		utils.RenderTemplate(w, http.StatusBadRequest, "profileEdit", pageData)
		return
	}
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

// render change password page
func (h *UseHandler) GetChangePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	pageData := models.PageData{
		User:        user,
		IsOwner:     true,
		PageContent: nil,
	}
	utils.RenderTemplate(w, http.StatusOK, "password_change", pageData)
}

func (h *UseHandler) PostChangePassword(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	if parseErr := r.ParseForm(); parseErr != nil {
		handleError(w, customerrors.ErrInternalError)
		return
	}

	currentPass := strings.TrimSpace(r.FormValue("currentPassword"))
	newPass := strings.TrimSpace(r.FormValue("newPassword"))
	confirmPass := strings.TrimSpace(r.FormValue("confirmPassword"))

	if newPass != confirmPass {
		pageData := models.PageData{
			User:    user,
			IsOwner: true,
			Error:   "password does not match confirmation password.",
		}
		utils.RenderTemplate(w, http.StatusBadRequest, "password_change", pageData)
		return
	}

	if UpdateErr := h.service.ChangePasswordService(cx, user.Id, currentPass, newPass); UpdateErr != nil {
		pageData := models.PageData{
			User:    user,
			IsOwner: true,
			Error:   UpdateErr.Error(),
		}
		utils.RenderTemplate(w, http.StatusBadRequest, "password_change", pageData)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// i dont think this it needed any more
// this does a server side extra check on the usename and email for duplicets and incorrect format
// func (h *UseHandler) CheckIfAvailabe(w http.ResponseWriter, r *http.Request) {
// 	cx := r.Context()
// 	var (
// 		exists   bool
// 		checkErr error
// 		message  string
// 	)
// 	switch r.URL.Path {
// 	case "/check-username":
// 		username := r.URL.Query().Get("username")
// 		// if !models.IsValidName(username) {
// 		// 	message = customerrors.ErrInvalidName.Error()
// 		// 	break
// 		// }
// 		exists, checkErr = h.service.CheckIfAvailable(cx, username)
// 		message = fmt.Sprintf("Username %s is not availabl", username)
// 	case "/check-email":
// 		email := r.URL.Query().Get("email")
// 		// if !models.IsValidEmail(email) {
// 		// 	message = customerrors.ErrInvalidData.Error()
// 		// 	break
// 		// }
// 		exists, checkErr = h.service.CheckIfAvailable(cx, email)
// 		message = customerrors.ErrDuplicateEmail.Error()
// 	default:
// 		handleError(w, customerrors.ErrInternalError)
// 		return
// 	}
// 	if checkErr != nil {
// 		http.Error(w, "Failed to Check Availability", http.StatusInternalServerError)
// 		return
// 	}

// 	payload := struct {
// 		Exists  bool   `json:"exists"`
// 		Message string `json:"message"`
// 	}{
// 		Exists:  exists,
// 		Message: message,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(payload)
// }
