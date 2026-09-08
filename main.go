package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitea.kood.tech/jyrkikarhunen/forum/database"
	"gitea.kood.tech/jyrkikarhunen/forum/handlers"
	"gitea.kood.tech/jyrkikarhunen/forum/middleware"
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
		middleware.Recoverer(handlers.HomePage)(w, r)
	})

	mux.HandleFunc("GET /user/profile", middleware.Recoverer(middleware.Restrict(sessionRepo, handlers.Profile)))

	mux.HandleFunc("GET /user/register", middleware.Recoverer(userHandler.GetRegisterUser))
	mux.HandleFunc("POST /user/register", middleware.Recoverer(userHandler.PostRegisterUser))

	mux.HandleFunc("GET /user/login", middleware.Recoverer(userHandler.GetSignInUser))
	mux.HandleFunc("POST /user/login", middleware.Recoverer(userHandler.SignInUser))

	mux.HandleFunc("POST /user/logout", middleware.Recoverer(userHandler.SignOutUser))

	mux.HandleFunc("GET /api/user/check-username", middleware.Recoverer(userHandler.CheckIfAvailabe))
	mux.HandleFunc("GET /api/user/check-email", middleware.Recoverer(userHandler.CheckIfAvailabe))

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
