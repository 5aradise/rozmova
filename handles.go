package main

import "net/http"

func createHandles(mux *http.ServeMux, cfg *apiConfig) {
	const filepathRoot = "."
	const appPathRoot = "/app"
	const apiPathRoot = "/api"
	const adminPathRoot = "/admin"

	mux.HandleFunc("GET "+apiPathRoot+"/healthz", healthz)
	mux.HandleFunc("GET "+adminPathRoot+"/metrics", cfg.showMetrics)
	mux.HandleFunc(""+apiPathRoot+"/reset", cfg.resetHits)

	mux.HandleFunc("GET "+apiPathRoot+"/messages", cfg.getMessages)
	mux.HandleFunc("GET "+apiPathRoot+"/messages/{messageId}", cfg.getMessageById)
	mux.HandleFunc("POST "+apiPathRoot+"/messages", cfg.postMessage)

	mux.HandleFunc("GET "+apiPathRoot+"/users", cfg.getUsers)
	mux.HandleFunc("GET "+apiPathRoot+"/users/{userId}", cfg.getUserById)
	mux.HandleFunc("POST "+apiPathRoot+"/users", cfg.postUser)

	mux.Handle(""+appPathRoot+"/*", cfg.middlewareMetricsInc(http.StripPrefix(appPathRoot, http.FileServer(http.Dir(filepathRoot)))))
}