package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/auth"
	"github.com/szmktk/chirpy/internal/database"
)

func (srv *Server) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	type response struct {
		User
	}

	ctxVal := r.Context().Value(contextKeyUserID)
	parsedUserID, ok := ctxVal.(uuid.UUID)
	if !ok {
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	decoder := json.NewDecoder(r.Body)
	payload := UserData{}
	err := decoder.Decode(&payload)
	if err != nil {
		srv.logger.Error("Error decoding JSON body: %s", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		srv.logger.Error("Error hashing user password: %s", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	user, err := srv.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             parsedUserID,
		Email:          payload.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		srv.logger.Error("Error updating user data: %s", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	return respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
