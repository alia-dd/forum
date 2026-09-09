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
		handleError(w, customerrors.ErrInternalError)
	}
	userData := models.UserLogin{
		Username: strings.TrimSpace(r.FormValue("username")),
		Password: strings.TrimSpace(r.FormValue("password")),
	}

	session, siginErr := h.service.AuthenticateUserService(cx, userData)
	if siginErr != nil {
		handleError(w, siginErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.SesssionId,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
		// this secures the cookie data but since we are using http
		// it will block the cookie it self so not usefull currenly
		// Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UseHandler) SignOutUser(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	cookie, cookieErr := r.Cookie("session_token")
	if cookieErr == nil {
		if DeleteSessionErr := h.service.LogoutService(cx, cookie.Value); DeleteSessionErr != nil {
			if !errors.Is(DeleteSessionErr, customerrors.ErrNotFound) {
				handleError(w, DeleteSessionErr)
				return
			}
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
