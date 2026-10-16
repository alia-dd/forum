package middleware

import (
	"context"
	"net/http"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/handlers"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func Recoverer(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				handlers.HandleError(w, r, customerrors.ErrInternalError)
			}
		}()
		handler(w, r)
	}
}

// for GET endpoint to view post and comment this will work
// it allows the session but its not necessary
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

		ctx := context.WithValue(r.Context(), "user_session", &user)
		handler(w, r.WithContext(ctx))

	}
}

// wrap this in request where the user must be authorized eg creating post/comment
func Restrict(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cx := r.Context()
		// we save the session token data when we create the session
		// it retrieves that and uses the cookie session uid to get the useid
		cookie, cookieErr := r.Cookie("session_token")
		if cookieErr != nil {
			utils.RedirectTologin(w, r)
			return
		}
		session, sessionErr := sessionRepo.GetSessionwithSessionId(cx, cookie.Value)
		if sessionErr != nil {
			utils.RedirectTologin(w, r)
			return
		}
		user, err := userRepo.FetchUserData(cx, session.UserID)
		if err != nil {
			utils.RedirectTologin(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user_session", &user)
		handler(w, r.WithContext(ctx))

	}
}

// wrap this in request where the user must be an admin
func Admin(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository, handler http.HandlerFunc) http.HandlerFunc {
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
		if err != nil || user.Role != "admin" {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "user_session", &user)
		handler(w, r.WithContext(ctx))
	}
}
