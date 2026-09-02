package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitea.kood.tech/jyrkikarhunen/forum/database"
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

	mux := http.NewServeMux()

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
