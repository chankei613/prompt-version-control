package main

import (
	"log"
	"net/http"

	"github.com/chankei613/prompt-version-control/internal/api"
	"github.com/chankei613/prompt-version-control/internal/db"
)

func main() {
	conn, err := db.Init("prompt-version-control.db")
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Println("prompt-version-control backend listening on :8425")
	if err := http.ListenAndServe(":8425", router); err != nil {
		log.Fatal(err)
	}
}
