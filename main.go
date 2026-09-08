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

	// GET / - fetches all posts. Optional, combinable query filters: ?category={id}  ?author={id|me|username}  ?liked=true
	mux.HandleFunc("GET /", postHandler.GetPosts)
	// GET /post/{id} - just an int
	mux.HandleFunc("GET /post/{id}", postHandler.GetPostByID)

	//below need auth
	// GET /post/new - loads template (once implemented) for submitting post
	mux.HandleFunc("GET /post/new", postHandler.NewPostForm)
	// POST /post/new - create post. Form: title, content, main_category={id}, category={id}… (sides)
	mux.HandleFunc("POST /post/new", postHandler.CreatePost)
	// GET /post/{id}/edit - loads template (once implemented) for editing a post you own
	mux.HandleFunc("GET /post/{id}/edit", postHandler.EditPostForm)
	// POST /post/{id}/edit - update your post. Form: title, content, main_category={id}, category={id}… (sides)
	mux.HandleFunc("POST /post/{id}/edit", postHandler.UpdatePost)

	// GET /category/new - loads template (once implemented) for creating new category
	mux.HandleFunc("GET /category/new", categoryHandler.NewCategoryForm)
	// POST /category/new - create category. Form: name
	mux.HandleFunc("POST /category/new", categoryHandler.CreateCategory)
	// GET /category/{id}/edit - loads template (once implemented) for editing existing category
	mux.HandleFunc("GET /category/{id}/edit", categoryHandler.EditCategoryForm)
	// POST /category/{id}/edit - rename category. Form: name
	mux.HandleFunc("POST /category/{id}/edit", categoryHandler.UpdateCategory)
	// POST /category/{id}/delete - delete category (if it has no references elsewhere)
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
