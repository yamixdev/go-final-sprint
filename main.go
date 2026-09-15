package main

import (
	"log"
	"os"

	"github.com/yamixdev/go-final-sprint/pkg/db"
	"github.com/yamixdev/go-final-sprint/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	err := db.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
