package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/database"
)

func (srv *Server) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	ordering := strings.ToLower(r.URL.Query().Get("sort"))
	authorID := r.URL.Query().Get("author_id")
	if authorID == "" {
		srv.logger.Info("Getting all chirps stored in the database")
		dbChirps, err := srv.db.GetAllChirps(r.Context())
		if err != nil {
			srv.logger.Error("Error getting all chirps: %s", "err", err)
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if ordering == "desc" {
			reverse(dbChirps)
			respondWithJSON(w, http.StatusOK, mapDbChirps(dbChirps))
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

	srv.logger.Info("Getting all chirps for user", "user_id", authorUUID)
	dbChirps, err := srv.db.GetAllChirpsForUser(r.Context(), authorUUID)
	if err != nil {
		srv.logger.Error("Error getting all user chirps: %s", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if ordering == "desc" {
		reverse(dbChirps)
		respondWithJSON(w, http.StatusOK, mapDbChirps(dbChirps))
		return
	}
	respondWithJSON(w, http.StatusOK, mapDbChirps(dbChirps))
}

func (srv *Server) HandlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing UUID: %s", err))
		return
	}

	dbChirp, err := srv.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		srv.logger.Info("Chirp not found", "err", err)
		respondWithError(w, http.StatusNotFound, "Chirp with given id has not been found")
		return
	}

	srv.logger.Info("Getting details of chirp", "id", chirpUUID)
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

func reverse(records []database.Chirp) {
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}
}
