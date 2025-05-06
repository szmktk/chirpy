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

func (srv *Server) CreateUser(w http.ResponseWriter, r *http.Request) error {
	type response struct {
		User
		Token string `json:"token,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := UserData{}
	err := decoder.Decode(&payload)
	if err != nil {
		return APIError{Status: http.StatusBadRequest, Msg: fmt.Sprintf("Error decoding JSON body: %s", err)}

	}

	if payload.Email == "" {
		return APIError{Status: http.StatusBadRequest, Msg: "Email cannot be empty"}

	}
	if payload.Password == "" {
		return APIError{Status: http.StatusBadRequest, Msg: "Password cannot be empty"}

	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		srv.logger.Error("Error hashing user password", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}

	}

	user, err := srv.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          payload.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return APIError{Status: http.StatusConflict, Msg: "A user with this email already exists"}
		} else {
			srv.logger.Error("Error creating user", "err", err)
			return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
		}

	}

	return respondWithJSON(w, http.StatusCreated, response{
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
