package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/JorgeToAn/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting API Key: %s", err)
		respondWithError(w, http.StatusUnauthorized, ErrorUnauthorized.String())
		return
	}
	if apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, ErrorUnauthorized.String())
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
		return
	}

	switch params.Event {
	case "user.upgraded":
		userID, err := uuid.Parse(params.Data.UserID)
		if err != nil {
			log.Printf("Error parsing UUID: %s", err)
			respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
			return
		}
		cfg.handleUserUpgradeWebhook(w, r, userID)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func (cfg *ApiConfig) handleUserUpgradeWebhook(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	err := cfg.DBQueries.UpgradeUserPlan(r.Context(), userID)
	if err != nil {
		log.Printf("Error updating user plan: %s", err)
		respondWithError(w, http.StatusInternalServerError, ErrorGeneric.String())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
