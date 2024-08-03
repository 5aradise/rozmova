package main

import "net/http"

func createHandles(mux *http.ServeMux, cfg *apiConfig) {
	const adminPathRoot = "/admin"
	const apiPathRoot = "/api"
	const appPathRoot = "/app"

	const messagesPath = "/messages"
	const usersPath = "/users"

	const filepathRoot = "./public"

	mux.HandleFunc("GET "+adminPathRoot+"/metrics", cfg.showMetrics)

	mux.HandleFunc("GET "+apiPathRoot+"/healthz", healthz)
	mux.HandleFunc(""+apiPathRoot+"/reset", cfg.resetHits)

	mux.HandleFunc("GET "+apiPathRoot+messagesPath, cfg.getMessages)
	mux.HandleFunc("GET "+apiPathRoot+messagesPath+"/{messageId}", cfg.getMessageById)
	mux.HandleFunc("POST "+apiPathRoot+messagesPath, cfg.postMessage)

	mux.HandleFunc("GET "+apiPathRoot+usersPath, cfg.getUsers)
	mux.HandleFunc("GET "+apiPathRoot+usersPath+"/{userId}", cfg.getUserById)
	mux.HandleFunc("POST "+apiPathRoot+usersPath, cfg.registerUser)
	mux.HandleFunc("PUT "+apiPathRoot+usersPath, cfg.changeUser)

	mux.HandleFunc("POST "+apiPathRoot+"/login", cfg.loginUser)
	mux.HandleFunc("POST "+apiPathRoot+"/refresh", cfg.refreshToken)
	mux.HandleFunc("POST "+apiPathRoot+"/revoke", cfg.revokeToken)

	mux.Handle(""+appPathRoot+"/*", cfg.middlewareMetricsInc(http.StripPrefix(appPathRoot, http.FileServer(http.Dir(filepathRoot)))))
}
