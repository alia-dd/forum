package handlers

import (
	"fmt"
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
		fmt.Println("user handler profile cookie session data >", user, ok)
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	fmt.Println(user)

	pageData := models.PageData{
		User:        user,
		IsOwner:     ok,
		PageContent: nil,
	}
	utils.RenderTemplate(w, http.StatusOK, "profile", pageData)
}
