package handlers

import (
	"net/http"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	pageData := models.PageData{
		User:        user,
		IsOwner:     ok,
		PageContent: nil,
	}

	utils.RenderTemplate(w, http.StatusOK, "home", pageData)
}

func Profile(w http.ResponseWriter, r *http.Request) {

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
	utils.RenderTemplate(w, http.StatusOK, "profile", pageData)
}
