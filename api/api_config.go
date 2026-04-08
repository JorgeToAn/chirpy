package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/JorgeToAn/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DBQueries      *database.Queries
	JWTSecret      string
	PolkaKey       string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) HandlerMetricsGet(w http.ResponseWriter, _ *http.Request) {
	template := `
	<html>
	  <body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	  </body>
	</html>`
	hits := cfg.FileserverHits.Load()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(template, hits)))
}

func (cfg *ApiConfig) HandlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		respondWithError(w, 403, "Forbidden action")
		return
	}

	err := cfg.DBQueries.DeleteUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
		respondWithError(w, 500, "Couldn't delete users")
		return
	}

	cfg.FileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset OK\n"))
}
