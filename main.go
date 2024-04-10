package main

import (
	"log"
	"net/http"

	"github.com/5aradise/rozmova/internal/database"
)

func main() {
	const port = "8080"
	const databasePath = "database.json"

	mux := http.NewServeMux()

	db, err := database.NewDB(databasePath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := NewApiConfig(db)

	createHandles(mux, cfg)

	corsMux := middlewareCors(mux)

	err = http.ListenAndServe(":"+port, corsMux)

	if err != nil {
		log.Fatal(err)
	}
}
