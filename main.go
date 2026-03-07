package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

// API Status
type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	port := "8080"
	filepathRoot := "."
	cfg := &apiConfig{}

	// Create a new ServeMux
	serverMux := http.NewServeMux()

	// Create a new HTTP server with the specified configuration
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}

	// Add handler
	serverMux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serverMux.HandleFunc("POST /api/reset", cfg.handlerReset)

	// Start the server
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",
		cfg.fileserverHits.Load())))
}
