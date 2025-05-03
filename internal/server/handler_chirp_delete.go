package server

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (srv *Server) DeleteChirp(w http.ResponseWriter, r *http.Request) error {
	ctxVal := r.Context().Value(contextKeyUserID)
	parsedUserID, ok := ctxVal.(uuid.UUID)
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		srv.logger.Warn("Error parsing UUID", "err", err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing UUID: %s", err))
	}

	dbChirp, err := srv.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		srv.logger.Info("Chirp not found", "err", err)
		return APIError{Status: http.StatusNotFound, Msg: "Chirp with given id has not been found"}
	}
	if dbChirp.UserID != parsedUserID {
		return APIError{Status: http.StatusForbidden, Msg: "Deleting chirps of other users is not allowed"}
	}

	if err := srv.db.DeleteChirp(r.Context(), chirpUUID); err != nil {
		srv.logger.Error("Error deleting chirp: %s", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	respondWithNoContent(w)
	return nil
}
