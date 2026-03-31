package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

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

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
	}

	mux := http.NewServeMux()

	// FILE SERVER
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// API
	mux.HandleFunc("GET /api/healthz", handlerHealth)

	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	// ADMIN
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricsGet)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Started server at port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
