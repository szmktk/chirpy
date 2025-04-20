package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	parsedUserID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		logger.Info("Error parsing UUID", "err", err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing UUID: %s", err))
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		logger.Info("Chirp not found", "err", err)
		respondWithError(w, http.StatusNotFound, "Chirp with given id has not been found")
		return
	}
	if dbChirp.UserID != parsedUserID {
		respondWithError(w, http.StatusForbidden, "Deleting chirps of other users is not allowed")
		return
	}

	if err := cfg.db.DeleteChirp(r.Context(), chirpUUID); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting chirp: %s", err))
		return
	}

	respondWithNoContent(w)
}
