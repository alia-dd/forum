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
	var deleteDB bool
	var seedDB bool

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-d" || os.Args[i] == "-delete" {
			deleteDB = true
		}
		if os.Args[i] == "-i" || os.Args[i] == "-init" {
			seedDB = true
		}
	}

	if deleteDB {
		err := os.Remove("data/forum.db")
		if err != nil {
			fmt.Println(err)
		}
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if seedDB {
		err := database.SeedData(db)
		if err != nil {
			fmt.Println(err)
		}
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
		middleware.Recoverer(middleware.AllowGuest(sessionRepo, userRepo, handlers.HomePage))(w, r)
	})

	mux.HandleFunc("GET /user/profile", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, handlers.Profile)))

	mux.HandleFunc("GET /user/register", middleware.Recoverer(userHandler.GetRegisterUser))
	mux.HandleFunc("POST /user/register", middleware.Recoverer(userHandler.PostRegisterUser))

	mux.HandleFunc("GET /user/login", middleware.Recoverer(userHandler.GetSignInUser))
	mux.HandleFunc("POST /user/login", middleware.Recoverer(userHandler.SignInUser))

	mux.HandleFunc("POST /user/logout", middleware.Recoverer(userHandler.SignOutUser))

	mux.HandleFunc("GET /api/user/check-username", middleware.Recoverer(userHandler.CheckIfAvailabe))
	mux.HandleFunc("GET /api/user/check-email", middleware.Recoverer(userHandler.CheckIfAvailabe))

	postRepo := repository.NewPostRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	postHandler := handlers.NewPostHandler(postRepo, categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)

	mux.HandleFunc("GET /", postHandler.GetPosts)
	mux.HandleFunc("GET /post/{id}", postHandler.GetPostByID)
	//below need auth
	mux.HandleFunc("GET /post/new", postHandler.NewPostForm)
	mux.HandleFunc("POST /post/new", postHandler.CreatePost)
	mux.HandleFunc("GET /post/{id}/edit", postHandler.EditPostForm)
	mux.HandleFunc("POST /post/{id}/edit", postHandler.UpdatePost)

	mux.HandleFunc("GET /category/new", categoryHandler.NewCategoryForm)
	mux.HandleFunc("POST /category/new", categoryHandler.CreateCategory)
	//above need auth

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
