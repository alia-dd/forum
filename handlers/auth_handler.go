package handlers

import (
	"errors"
	"net/http"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func (h *UseHandler) GetSignInUser(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, http.StatusOK, "login_user", nil)
}

func (h *UseHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()

	if parseErr := r.ParseForm(); parseErr != nil {
		handleError(w, customerrors.ErrInternalError, "", 0)
	}
	userData := models.UserLogin{
		Username: strings.TrimSpace(r.FormValue("username")),
		Password: strings.TrimSpace(r.FormValue("password")),
	}

	session, siginErr := h.service.AuthenticateUserService(cx, userData)
	if siginErr != nil {
		handleError(w, siginErr, "", 0)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.SesssionId,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (h *UseHandler) SignOutUser(w http.ResponseWriter, r *http.Request) {
	if DeleteSessionErr := h.service.DeleteUserSession(); DeleteSessionErr != nil {
		if !errors.Is(DeleteSessionErr, customerrors.ErrNotFound) {
			handleError(w, DeleteSessionErr, "", 0)
		}
	}
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
