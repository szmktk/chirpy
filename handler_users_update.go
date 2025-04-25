package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/auth"
	"github.com/szmktk/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
	}

	ctxVal := r.Context().Value(contextKeyUserID)
	parsedUserID, ok := ctxVal.(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	payload := UserData{}
	err := decoder.Decode(&payload)
	if err != nil {
		logger.Error("Error decoding JSON body: %s", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		logger.Error("Error hashing user password: %s", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             parsedUserID,
		Email:          payload.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		logger.Error("Error updating user data: %s", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
