package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

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
	var dbChirps []database.Chirp

	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Author ID is invalid")
		return
	}

	sorting, err := sortingFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Sort is invalid")
		return
	}

	if authorID != uuid.Nil {
		dbChirps, err = cfg.DBQueries.GetChirpsByAuthor(r.Context(), authorID)
	} else {
		dbChirps, err = cfg.DBQueries.GetChirps(r.Context())
	}
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
		return
	}

	if sorting == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool {
			return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt)
		})
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

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}

	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return authorID, nil
}

func sortingFromRequest(r *http.Request) (string, error) {
	sorting := r.URL.Query().Get("sort")
	if sorting == "" {
		return "asc", nil
	}
	if sorting != "asc" && sorting != "desc" {
		return "", fmt.Errorf("invalid sorting value: %s", sorting)
	}
	return sorting, nil
}
