package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/database"
)

func (srv *Server) GetAllChirps(w http.ResponseWriter, r *http.Request) error {
	ordering := strings.ToLower(r.URL.Query().Get("sort"))
	authorID := r.URL.Query().Get("author_id")
	if authorID == "" {
		srv.logger.Info("Getting all chirps stored in the database")
		dbChirps, err := srv.db.GetAllChirps(r.Context())
		if err != nil {
			srv.logger.Error("Error getting all chirps: %s", "err", err)
			return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
		}
		if ordering == "desc" {
			reverse(dbChirps)
			return respondWithJSON(w, http.StatusOK, mapDbChirps(dbChirps))
		}

		return respondWithJSON(w, http.StatusOK, mapDbChirps(dbChirps))

	}

	authorUUID, err := uuid.Parse(authorID)
	if err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: fmt.Sprintf("Error parsing UUID: %s", err)}
	}

	srv.logger.Info("Getting all chirps for user", "user_id", authorUUID)
	dbChirps, err := srv.db.GetAllChirpsForUser(r.Context(), authorUUID)
	if err != nil {
		srv.logger.Error("Error getting all user chirps: %s", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	if ordering == "desc" {
		reverse(dbChirps)
		return respondWithJSON(w, http.StatusOK, mapDbChirps(dbChirps))

	}
	return respondWithJSON(w, http.StatusOK, mapDbChirps(dbChirps))
}

func (srv *Server) GetChirp(w http.ResponseWriter, r *http.Request) error {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: fmt.Sprintf("Error parsing UUID: %s", err)}
	}

	dbChirp, err := srv.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		srv.logger.Info("Chirp not found", "err", err)
		return APIError{Status: http.StatusNotFound, Msg: "Chirp with given id has not been found"}
	}

	srv.logger.Info("Getting details of chirp", "id", chirpUUID)
	return respondWithJSON(w, http.StatusOK, Chirp{
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
