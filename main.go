package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitea.kood.tech/jyrkikarhunen/forum/database"
	"gitea.kood.tech/jyrkikarhunen/forum/handlers"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
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

	mux := http.NewServeMux()

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
	mux.HandleFunc("GET /category/{id}/edit", categoryHandler.EditCategoryForm)
	mux.HandleFunc("POST /category/{id}/edit", categoryHandler.UpdateCategory)
	mux.HandleFunc("POST /category/{id}/delete", categoryHandler.DeleteCategory)
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
