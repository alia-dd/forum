package handlers

import (
	"net/http"

	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, http.StatusOK, "home", nil)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	// _, ok := r.Context().Value("session_key").(*models.Session)
	// if !ok {
	// 	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	// }
	utils.RenderTemplate(w, http.StatusOK, "profile", nil)
}
