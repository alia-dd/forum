package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitea.kood.tech/jyrkikarhunen/forum/handlers"
	"gitea.kood.tech/jyrkikarhunen/forum/middleware"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/service"
)

func InitServer(db *sql.DB) {
	mux := http.NewServeMux()

	config := initRepoHandlers(db)

	registerRoutes(mux, config)

	runServer(mux)
}

func runServer(mux *http.ServeMux) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	fmt.Println("Server starting on http://localhost:8080")

	log.Fatal(server.ListenAndServe())
}

func initRepoHandlers(db *sql.DB) routeDependencies {
	d := routeDependencies{}

	sessionRepo := repository.NewSessionRepository(db)
	userRepo := repository.NewUserRepository(db)
	reactRepo := repository.NewReactionRepository(db)
	postRepo := repository.NewPostRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	searchRepo := repository.NewSearchRepository(db)

	userService := service.NewUserService(userRepo, sessionRepo)
	reactService := service.NewReactionService(reactRepo)

	d.userRepo = userRepo
	d.sessionRepo = sessionRepo
	d.searchHandler = handlers.NewSearchHandler(searchRepo)
	d.reactHandler = handlers.NewReactionHandler(reactService, postRepo, commentRepo)
	d.userHandler = handlers.NewUserHandler(userService, postRepo, commentRepo, categoryRepo, *searchRepo)
	d.postHandler = handlers.NewPostHandler(postRepo, categoryRepo, commentRepo)
	d.categoryHandler = handlers.NewCategoryHandler(categoryRepo)
	d.commentHandler = handlers.NewCommentHandler(commentRepo)

	return d
}

// middleware wrappers
func handle(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.HandleFunc(pattern, middleware.Recoverer(handler))
}

func handleRestricted(mux *http.ServeMux, pattern string, d routeDependencies, handler http.HandlerFunc) {
	handle(mux, pattern, middleware.Restrict(d.sessionRepo, d.userRepo, handler))
}

func handleGuest(mux *http.ServeMux, pattern string, d routeDependencies, handler http.HandlerFunc) {
	handle(mux, pattern, middleware.AllowGuest(d.sessionRepo, d.userRepo, handler))
}

func handleAdmin(mux *http.ServeMux, pattern string, d routeDependencies, handler http.HandlerFunc) {
	handle(mux, pattern, middleware.Admin(d.sessionRepo, d.userRepo, handler))
}
