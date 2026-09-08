package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitea.kood.tech/jyrkikarhunen/forum/database"
	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/handlers"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/service"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func main() {
	var seed bool

	if len(os.Args) == 2 {
		if os.Args[1] == "-i" || os.Args[1] == "-init" {
			seed = true
		}
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if seed {
		database.SeedData(db)
	}

	utils.InitializeTemplate()

	mux := http.NewServeMux()

	sessionRepo := repository.NewSessionRepository(db)
	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo, sessionRepo)
	userHandler := handlers.NewUserHandler(userService)

	mux.Handle("GET /static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			utils.RenderTemplate(w, http.StatusNotFound, "error", &models.ErrorStruct{Error: "400", ErrorMs: "Page Not Found"})
			return
		}
		MiddleWare(db, handlers.HomePage)(w, r)
	})

	// mux.HandleFunc("GET /", MiddleWare(db, handlers.HomePage))

	mux.HandleFunc("GET /user/register", MiddleWare(db, userHandler.GetRegisterUser))
	mux.HandleFunc("POST /user/register", MiddleWare(db, userHandler.PostRegisterUser))

	mux.HandleFunc("GET /user/login", MiddleWare(db, userHandler.GetSignInUser))
	mux.HandleFunc("POST /user/login", MiddleWare(db, userHandler.SignInUser))

	mux.HandleFunc("POST /user/logout", MiddleWare(db, userHandler.SignOutUser))

	mux.HandleFunc("GET /api/user/check-username", MiddleWare(db, userHandler.CheckIfAvailabe))
	mux.HandleFunc("GET /api/user/check-email", MiddleWare(db, userHandler.CheckIfAvailabe))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	fmt.Println("Server starting on :8080")

	log.Fatal(server.ListenAndServe())
}

func MiddleWare(db *sql.DB, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				http.Error(w, customerrors.ErrInternalError.Error(), http.StatusInternalServerError)
			}

		}()
		handler(w, r)

	}
}
