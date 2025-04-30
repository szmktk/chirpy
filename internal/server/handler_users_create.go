package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/auth"
	"github.com/szmktk/chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (srv *Server) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
		Token string `json:"token,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := UserData{}
	err := decoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding JSON body: %s", err))
		return
	}

	if payload.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email cannot be empty")
		return
	}
	if payload.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password cannot be empty")
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		srv.logger.Error("Error hashing user password: %s", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user, err := srv.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          payload.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			respondWithError(w, http.StatusConflict, "A user with this email already exists")
		} else {
			srv.logger.Error("Error creating user: %s", "err", err)
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}

// isDuplicateKeyError checks if the error is a result of a duplicate key constraint.
func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value")
}
