package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *sql.DB
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbConn,
		platform:       platform,
	}

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerChirpsGet)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
