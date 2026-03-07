package main

import (
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

	//Add handlers
	fshandler := cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	serverMux.Handle("/app/", fshandler)

	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	serverMux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	// Create a new HTTP server with the specified configuration
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}

	// Start the server
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
