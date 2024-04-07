package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const appPathRoot = "/app"
	const apiPathRoot = "/api"
	const port = "8080"

	mux := http.NewServeMux()

	cfg := NewApiConfig()

	mux.HandleFunc("GET "+apiPathRoot+"/healthz", healthz)
	mux.HandleFunc("GET "+apiPathRoot+"/metrics", cfg.showMetrics)
	mux.HandleFunc(""+apiPathRoot+"/reset", cfg.resetHits)
	mux.Handle(""+appPathRoot+"/*", cfg.middlewareMetricsInc(http.StripPrefix(appPathRoot, http.FileServer(http.Dir(filepathRoot)))))

	corsMux := middlewareCors(mux)

	err := http.ListenAndServe(":"+port, corsMux)

	if err != nil {
		log.Fatal(err)
	}
}
