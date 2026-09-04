package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitea.kood.tech/jyrkikarhunen/forum/database"
	"gitea.kood.tech/jyrkikarhunen/forum/handlers"
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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	mux.HandleFunc("POST /api/user/register", MiddleWare(db, userHandler.CreateUser))
	mux.HandleFunc("POST /api/user/{id}", MiddleWare(db, userHandler.CreateUser))
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
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

		}()
		handler(w, r)

	}
}
