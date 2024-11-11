package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const maxChirpLength = 140

var forbiddenWords = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

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

	cleanedBody := sanitizeInput(payload.Body)

	respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": cleanedBody})
}

func sanitizeInput(body string) string {
	var sanitized []string
	for _, word := range strings.Split(body, " ") {
		if forbiddenWords[strings.ToLower(word)] {
			word = "****"
		}
		sanitized = append(sanitized, word)
	}
	return strings.Join(sanitized, " ")
}
