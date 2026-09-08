package handlers

import (
	"fmt"
	"net/http"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	var pageData models.PageData
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	fmt.Println(ok)
	pageData = models.PageData{
		User:        user,
		PageContent: nil,
	}

	fmt.Println(pageData)
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
		PageContent: nil,
	}
	utils.RenderTemplate(w, http.StatusOK, "profile", pageData)
}
