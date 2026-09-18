package handlers

import (
	"fmt"
	"net/http"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	pageData := models.PageData[MainPage]{
		User:        user,
		IsOwner:     ok,
		PageContent: MainPage{},
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

	// replace mainpage with a dedicated profile struct later
	pageData := models.PageData[MainPage]{
		User:        user,
		IsOwner:     ok,
		PageContent: MainPage{},
	}
	utils.RenderTemplate(w, http.StatusOK, "profile", pageData)
}
