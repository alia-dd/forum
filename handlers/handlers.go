package handlers

import (
	"net/http"

	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, http.StatusOK, "home", nil)
}
