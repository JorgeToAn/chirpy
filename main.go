package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/JorgeToAn/chirpy/api"
	"github.com/JorgeToAn/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = "8080"
const filepathRoot = "."

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("missing JWT_SECRET environment variable")
	}

	apiCfg := api.ApiConfig{
		FileserverHits: atomic.Int32{},
		DBQueries:      database.New(db),
		JWTSecret:      jwtSecret,
	}

	mux := http.NewServeMux()

	// FILE SERVER
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// API
	mux.HandleFunc("GET /api/healthz", api.HandlerHealth)

	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.HandlerGetChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerGetAllChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerCreateChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandlerDeleteChirp)

	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRevoke)

	mux.HandleFunc("POST /api/users", apiCfg.HandlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUpdateUser)

	// ADMIN
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetricsGet)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerMetricsReset)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Started server at port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
