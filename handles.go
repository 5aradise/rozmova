package main

import "net/http"

func createHandles(mux *http.ServeMux, cfg *apiConfig) {
	const filepathRoot = "."
	
	const appPathRoot = "/app"
	const apiPathRoot = "/api"
	const adminPathRoot = "/admin"

	const messagesPath = "/messages"
	const usersPath = "/users"
	const loginPath = "/login"

	mux.HandleFunc("GET "+apiPathRoot+"/healthz", healthz)
	mux.HandleFunc("GET "+adminPathRoot+"/metrics", cfg.showMetrics)
	mux.HandleFunc(""+apiPathRoot+"/reset", cfg.resetHits)

	mux.HandleFunc("GET "+apiPathRoot+messagesPath, cfg.getMessages)
	mux.HandleFunc("GET "+apiPathRoot+messagesPath+"/{messageId}", cfg.getMessageById)
	mux.HandleFunc("POST "+apiPathRoot+messagesPath, cfg.postMessage)

	mux.HandleFunc("GET "+apiPathRoot+usersPath, cfg.getUsers)
	mux.HandleFunc("GET "+apiPathRoot+usersPath+"/{userId}", cfg.getUserById)
	mux.HandleFunc("POST "+apiPathRoot+usersPath, cfg.registerUser)

	mux.HandleFunc("POST "+apiPathRoot+loginPath, cfg.loginUser)

	mux.Handle(""+appPathRoot+"/*", cfg.middlewareMetricsInc(http.StripPrefix(appPathRoot, http.FileServer(http.Dir(filepathRoot)))))
}
