package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/JorgeToAn/chirpy/internal/auth"
	"github.com/JorgeToAn/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, ErrorUnauthorized.String())
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, ErrorUnauthorized.String())
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	const maxLength int = 140

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
		return
	}

	if len(params.Body) > maxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	dbChirp, err := cfg.DBQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   replaceWords(params.Body, "****", profaneWords),
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *ApiConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DBQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
		chirps = append(chirps, chirp)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *ApiConfig) HandlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Printf("Error parsing chirp ID: %s", err)
		respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
		return
	}

	dbChirp, err := cfg.DBQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *ApiConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, ErrorUnauthorized.String())
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("Error validating access token: %s", err)
		respondWithError(w, http.StatusUnauthorized, ErrorUnauthorized.String())
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		log.Printf("Error parsing UUID: %s", err)
		respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
		return
	}

	dbChirp, err := cfg.DBQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}

	if userID != dbChirp.UserID {
		respondWithError(w, http.StatusForbidden, ErrorForbidden.String())
		return
	}

	err = cfg.DBQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
