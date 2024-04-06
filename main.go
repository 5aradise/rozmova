package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const addrespathRoot = "/app"
	const port = "8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthz)
	mux.Handle(addrespathRoot+"/*", http.StripPrefix(addrespathRoot, http.FileServer(http.Dir(filepathRoot))))

	corsMux := middlewareCors(mux)

	err := http.ListenAndServe(":"+port, corsMux)

	if err != nil {
		log.Fatal(err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Status: OK`))
}
