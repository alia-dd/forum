package middleware

import (
	"context"
	"net/http"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
)

func Recoverer(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, customerrors.ErrInternalError.Error(), http.StatusInternalServerError)
			}
		}()
		handler(w, r)
	}
}

func AllowGuest(sessionRepo *repository.SessionRepository, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cx := r.Context()

		cookie, cookieErr := r.Cookie("session_token")
		if cookieErr != nil {
			handler(w, r)
			return
		}
		session, sessionErr := sessionRepo.GetSessionwithSessionId(cx, cookie.Value)
		if sessionErr != nil {
			handler(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "session_key", session)
		handler(w, r.WithContext(ctx))

	}
}

func Restrict(sessionRepo *repository.SessionRepository, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cx := r.Context()
		cookie, cookieErr := r.Cookie("session_token")
		if cookieErr != nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		session, sessionErr := sessionRepo.GetSessionwithSessionId(cx, cookie.Value)
		if sessionErr != nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), "session_key", session)
		handler(w, r.WithContext(ctx))

	}
}
