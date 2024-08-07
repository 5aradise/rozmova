package main

import (
	"log"
	"net/http"
	"os"

	"github.com/5aradise/rozmova/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	databasePath := os.Getenv("DATABSE_PATH")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	if port == "" {
		port = "8080"
	}
	if databasePath == "" {
		databasePath = "db"
	}
	if jwtSecret == "" {
		jwtSecret = "Glory to Ukraine!"
	}

	mux := http.NewServeMux()

	db, err := database.NewDB(databasePath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := NewApiConfig(db, jwtSecret, polkaKey)

	createHandles(mux, cfg)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Fatal(srv.ListenAndServe())
}
