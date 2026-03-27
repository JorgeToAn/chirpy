package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type response struct {
		Error string `json:"error"`
	}

	resp := response{
		Error: msg,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	const maxLength int = 140

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > maxLength {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type responsePayload struct {
		Valid bool `json:"valid"`
	}
	payload := responsePayload{
		Valid: true,
	}
	respondWithJSON(w, 200, payload)
}
