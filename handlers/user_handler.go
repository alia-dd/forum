package handlers

import (
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

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	r.ParseMultipartForm(10 << 20)

	// change the handleErr to return to htmx page instead of new page
	if parseErr := r.ParseForm(); parseErr != nil {
		handleError(w, r, customerrors.ErrInternalError)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	name := strings.TrimSpace(r.FormValue("fullname"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))
	conformPassword := strings.TrimSpace(r.FormValue("confirm_password"))

	var fileName string

	file, header, fileErr := r.FormFile("imagepath")
	if fileErr == nil && header.Size > 0 {
		defer file.Close()
		savename, saveErr := SaveUploadedFile(file)
		if saveErr != nil {
			pageData := models.PageData[models.UserRegister]{
				User:    nil,
				IsOwner: false,
				PageContent: models.UserRegister{
					Username:        username,
					Name:            name,
					Email:           email,
					Password:        password,
					ConformPassword: conformPassword,
				},
				Error: customerrors.ErrBadRequest.Error(),
			}
			utils.RenderTemplate(w, http.StatusBadRequest, "user_registration", pageData)
			return
		}
		fileName = savename
	}

	userData := models.UserRegister{
		Username:        username,
		Name:            name,
		Email:           email,
		Image:           fileName,
		Password:        password,
		ConformPassword: conformPassword,
	}

	if PostErr := h.service.CreateUserService(cx, userData); PostErr != nil {
		pageData := models.PageData[models.UserRegister]{
			User:    nil,
			IsOwner: false,
			PageContent: models.UserRegister{
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
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	// replace mainpage with a dedicated profile struct later
	pageData := models.PageData[MainPage]{
		User:        user,
		IsOwner:     ok,
		PageContent: MainPage{},
	}

	utils.RenderTemplate(w, http.StatusOK, "profileEdit", pageData)
}

// shows profile data for other registred users
func (h *UseHandler) GetOtherUserProfile(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	username := r.PathValue("username")

	profileData, fetchErr := h.service.GetUserService(cx, username)
	if fetchErr != nil {
		handleError(w, r, fetchErr)
		return
	}

	pageData := models.PageData[models.PublicUserInfo]{
		User:        user,
		IsOwner:     ok && user.ID == profileData.ID,
		PageContent: profileData,
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

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	r.ParseMultipartForm(10 << 20)

	if parseErr := r.ParseForm(); parseErr != nil {
		handleError(w, r, customerrors.ErrInternalError)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	name := strings.TrimSpace(r.FormValue("fullname"))
	email := strings.TrimSpace(r.FormValue("email"))
	bio := strings.TrimSpace(r.FormValue("bio"))
	var fileName string

	file, header, fileErr := r.FormFile("imagepath")
	if fileErr == nil && header.Size > 0 {
		defer file.Close()
		savename, saveErr := SaveUploadedFile(file)
		if saveErr != nil {
			// replace mainpage with a dedicated profile struct later
			pageData := models.PageData[MainPage]{
				User: &models.UserInfo{
					Username: username,
					Name:     name,
					Email:    email,
					Bio:      bio,
					Image:    fileName,
				},
				IsOwner:     ok,
				PageContent: MainPage{},
				Error:       customerrors.ErrBadRequest.Error(),
			}
			utils.RenderTemplate(w, http.StatusBadRequest, "profileEdit", pageData)
			return
		}
		fileName = savename
	}

	userData := models.UserUpdate{
		ID:  user.ID,
		Bio: &bio,
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
	if fileName != "" {
		userData.Image = &fileName
	}
	if UpdateErr := h.service.UpdateUserService(cx, userData); UpdateErr != nil {
		// replace mainpage with a dedicated profile struct later
		pageData := models.PageData[MainPage]{
			User: &models.UserInfo{
				ID:       user.ID,
				Username: username,
				Name:     name,
				Email:    email,
				Bio:      bio,
				Image:    fileName,
			},
			IsOwner:     ok,
			PageContent: MainPage{},
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
	// replace mainpage with a dedicated profile struct later
	pageData := models.PageData[MainPage]{
		User:        user,
		IsOwner:     true,
		PageContent: MainPage{},
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
		handleError(w, r, customerrors.ErrInternalError)
		return
	}

	currentPass := strings.TrimSpace(r.FormValue("currentPassword"))
	newPass := strings.TrimSpace(r.FormValue("newPassword"))
	confirmPass := strings.TrimSpace(r.FormValue("confirmPassword"))

	if newPass != confirmPass {
		// replace mainpage with a dedicated profile struct later
		pageData := models.PageData[MainPage]{
			User:    user,
			IsOwner: true,
			Error:   "password does not match confirmation password.",
		}
		utils.RenderTemplate(w, http.StatusBadRequest, "password_change", pageData)
		return
	}

	if UpdateErr := h.service.ChangePasswordService(cx, user.ID, currentPass, newPass); UpdateErr != nil {
		// replace mainpage with a dedicated profile struct later
		pageData := models.PageData[MainPage]{
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
