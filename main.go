package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/a6u6akera6/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// API Status
type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not connet to SQL DB")
	}

	dbQueries := database.New(db)

	port := "8080"
	filepathRoot := "."
	cfg := &apiConfig{
		db: dbQueries,
	}

	// Create a new ServeMux
	serverMux := http.NewServeMux()

	//Add handlers
	fshandler := cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	serverMux.Handle("/app/", fshandler)

	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	//serverMux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	serverMux.HandleFunc("POST /api/users", cfg.handlerUsers)
	serverMux.HandleFunc("POST /api/chirps", cfg.handlerChirps)

	// Create a new HTTP server with the specified configuration
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}

	// Start the server
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
