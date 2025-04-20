package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	if authorID == "" {
		logger.Info("Getting all chirps stored in the database")
		dbChirps, err := cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting all chirps: %s", err))
			return
		}

		respondWithJSON(w, http.StatusOK, mapDbChirps(dbChirps))
		return
	}

	authorUUID, err := uuid.Parse(authorID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing UUID: %s", err))
		return
	}

	logger.Info("Getting all chirps for user", "user_id", authorUUID)
	dbChirps, err := cfg.db.GetAllChirpsForUser(r.Context(), authorUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting all user chirps: %s", err))
		return
	}
	respondWithJSON(w, http.StatusOK, mapDbChirps(dbChirps))
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing UUID: %s", err))
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		logger.Info("Chirp not found", "err", err)
		respondWithError(w, http.StatusNotFound, "Chirp with given id has not been found")
		return
	}

	logger.Info("Getting details of chirp", "id", chirpUUID)
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

func mapDbChirps(entities []database.Chirp) []Chirp {
	chirps := make([]Chirp, 0)
	for _, c := range entities {
		chirps = append(chirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
	return chirps
}
