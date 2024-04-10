package main

import (
	"fmt"
	"net/http"

	"github.com/5aradise/rozmova/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func NewApiConfig(db *database.DB) *apiConfig {
	return &apiConfig{
		fileserverHits: 0,
		db:             db,
	}
}

func (cfg *apiConfig) middlewareMetricsInc(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		h.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}

func (cfg *apiConfig) showMetrics(w http.ResponseWriter, r *http.Request) {
	htmlMetrics := fmt.Sprintf(`
	<html>

	<body>
			<h1>Welcome, Rozmova Admin</h1>
			<p>Chirpy has been visited %d times!</p>
	</body>
	
	</html>
	`, cfg.fileserverHits)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlMetrics))
}
