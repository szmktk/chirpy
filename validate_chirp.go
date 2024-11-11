package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const maxChirpLength = 140

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := input{}
	err := decoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding JSON body: %s", err))
		return
	}

	if len(payload.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]bool{"valid": true})
}
