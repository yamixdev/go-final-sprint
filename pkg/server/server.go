package server

import (
	"log"
	"net/http"
	"os"

	"github.com/yamixdev/go-final-sprint/pkg/api"
)

func Run() {
	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540"
	}

	webDir := "./web"

	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	api.Init()

	log.Printf("Server started at http://localhost:%s", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
