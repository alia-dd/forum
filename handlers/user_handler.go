package handlers

import (
	"net/http"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/service"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

type UseHandler struct {
	service     *service.UserService
	postRep     repository.PostRepository
	commentRep  repository.CommentRepository
	categoryRep repository.CategoryRepository
	searchRep   repository.SearchRepository
}

func NewUserHandler(service *service.UserService, postRep repository.PostRepository, commentRep repository.CommentRepository, categoryRep repository.CategoryRepository, searchRep repository.SearchRepository) *UseHandler {
	return &UseHandler{service: service, postRep: postRep, commentRep: commentRep, categoryRep: categoryRep, searchRep: searchRep}
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

	utils.RedirectTologin(w, r)
}

func (h *UseHandler) GetEditUserProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
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

func (h *UseHandler) showProfile(w http.ResponseWriter, r *http.Request, username string) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	tab := r.URL.Query().Get("tab")
	if tab != "posts" && tab != "comments" && tab != "likedPosts" {
		tab = "posts"
	}
	pageContent := models.ProfilePage{Tab: tab}
	profileData, fetchErr := h.service.GetUserService(cx, username)
	if fetchErr != nil {
		handleError(w, r, fetchErr)
		return
	}
	userID := -1
	if ok {
		userID = user.ID
	}
	pageContent.Profile = profileData

	if tab == "posts" {
		posts, fetchPostErr := h.postRep.GetPost(cx, models.PostFilter{AuthorID: &profileData.ID}, userID)
		if fetchPostErr != nil {
			handleError(w, r, fetchPostErr)
			return
		}

		pageContent.Posts = posts
	}

	if tab == "comments" {
		comments, fetchCommErr := h.commentRep.GetCommentsByUser(cx, profileData.ID)
		if fetchCommErr != nil {
			handleError(w, r, fetchCommErr)
			return
		}
		pageContent.Comments = comments
	}

	if tab == "likedPosts" {
		likedPosts, fetchPostLikedErr := h.postRep.GetPost(cx, models.PostFilter{LikedByID: &profileData.ID}, userID)
		if fetchPostLikedErr != nil {
			handleError(w, r, fetchPostLikedErr)
			return
		}
		pageContent.LikedPosts = likedPosts
	}
	pageData := models.PageData[models.ProfilePage]{
		User:        user,
		IsOwner:     ok && user.ID == profileData.ID,
		PageContent: pageContent,
	}
	utils.RenderTemplate(w, http.StatusOK, "profile", pageData)

}

// shows profile data for other registred users
func (h *UseHandler) GetOtherUserProfile(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	h.showProfile(w, r, username)
}

func (h *UseHandler) GetOwnProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}
	h.showProfile(w, r, user.Username)
}

func (h *UseHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
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

// delete user using the corrent cookie userid
func (h *UseHandler) DeleteUserAcccount(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}
	if deleteErr := h.service.DeleteUser(cx, user.ID); deleteErr != nil {
		handleError(w, r, deleteErr)
		return
	}
	cookie, cookieErr := r.Cookie("session_token")
	if cookieErr == nil {
		h.service.LogoutService(cx, cookie.Value)
	}
	utils.RedirectTologin(w, r)
}

// render change password page
func (h *UseHandler) GetChangePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
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
		utils.RedirectTologin(w, r)
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
	utils.RedirectTologin(w, r)
}
