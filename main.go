package main

import (
	"log"
	"net/http"
)

func main() {

	port := "8080"
	filepathRoot := "."

	// Create a new ServeMux
	serverMux := http.NewServeMux()

	// Create a new HTTP server with the specified configuration
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}

	// Add handler
	serverMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	serverMux.HandleFunc("/healthz", handlerReadiness)

	// Start the server
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
