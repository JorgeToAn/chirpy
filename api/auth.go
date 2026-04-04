package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JorgeToAn/chirpy/internal/auth"
	"github.com/JorgeToAn/chirpy/internal/database"
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	dbUser, err := cfg.DBQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		log.Printf("Error checking password hash: %s", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	if !match {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	expiresIn := time.Duration(params.ExpiresInSeconds) * time.Second
	if expiresIn.Seconds() == 0 || expiresIn > defaultTokenExpiration {
		expiresIn = defaultTokenExpiration
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.JWTSecret, expiresIn)
	if err != nil {
		log.Printf("Error making JWT: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	dbRefreshToken, err := cfg.DBQueries.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:     auth.MakeRefreshToken(),
			ExpiresAt: time.Now().Add(refreshTokenExpiration).UTC(),
			UserID:    dbUser.ID,
		},
	)
	if err != nil {
		log.Printf("Error making refresh token: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	user := LoginResponse{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
		Token:        token,
		RefreshToken: dbRefreshToken.Token,
	}

	respondWithJSON(w, 200, user)
}
