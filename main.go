package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Readiness endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Fileserver endpoint
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	log.Println("Server listening on :8080")
	log.Fatal(server.ListenAndServe())
}
