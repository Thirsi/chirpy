package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hits: %d\n", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hits reset to 0\n")
}

func main() {
	mux := http.NewServeMux()

	apiCfg := &apiConfig{}

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Readiness endpoint - GET only
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Fileserver endpoint with middleware
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServer))

	// Metrics endpoint - GET only
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)

	// Reset endpoint - POST only
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)

	log.Println("Server listening on :8080")
	log.Fatal(server.ListenAndServe())
}
