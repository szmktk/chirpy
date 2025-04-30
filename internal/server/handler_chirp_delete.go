package server

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (srv *Server) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	ctxVal := r.Context().Value(contextKeyUserID)
	parsedUserID, ok := ctxVal.(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		srv.logger.Warn("Error parsing UUID", "err", err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing UUID: %s", err))
		return
	}

	dbChirp, err := srv.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		srv.logger.Info("Chirp not found", "err", err)
		respondWithError(w, http.StatusNotFound, "Chirp with given id has not been found")
		return
	}
	if dbChirp.UserID != parsedUserID {
		respondWithError(w, http.StatusForbidden, "Deleting chirps of other users is not allowed")
		return
	}

	if err := srv.db.DeleteChirp(r.Context(), chirpUUID); err != nil {
		srv.logger.Error("Error deleting chirp: %s", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithNoContent(w)
}
