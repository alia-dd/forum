package middleware

import (
	"context"
	"fmt"
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

// for GET endpoint to view post and comment this will work
// it allows the session but its not nessecery
func AllowGuest(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository, handler http.HandlerFunc) http.HandlerFunc {
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
		user, err := userRepo.FetchUserData(cx, session.UserID)
		if err != nil {
			handler(w, r)
			return
		}
		fmt.Println("middle", user)

		ctx := context.WithValue(r.Context(), "user_session", &user)
		handler(w, r.WithContext(ctx))

	}
}

// wrapp this in request where the user must be authorized eg creating post/comment
func Restrict(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository, handler http.HandlerFunc) http.HandlerFunc {
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
		user, err := userRepo.FetchUserData(cx, session.UserID)
		if err != nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "user_session", &user)
		handler(w, r.WithContext(ctx))

	}
}
