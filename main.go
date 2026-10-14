package main

import (
	"log"
	"os"

	"gitea.kood.tech/jyrkikarhunen/forum/database"
	"gitea.kood.tech/jyrkikarhunen/forum/server"
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
			log.Fatal(err)
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
			log.Fatal(err)
		}
	}

	utils.InitializeTemplate()

	server.InitServer(db)
}
