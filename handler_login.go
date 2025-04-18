package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/szmktk/chirpy/internal/auth"
)

const defaultTokenExpirySeconds int = 3600

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		TokenExpirySeconds *int   `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := input{}
	tokenExpirySeconds := defaultTokenExpirySeconds
	err := decoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding JSON body: %s", err))
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
	if payload.TokenExpirySeconds != nil {
		tokenExpirySeconds = *payload.TokenExpirySeconds
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		logger.Info("User not found", "err", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if err := auth.CheckPasswordHash(payload.Password, user.HashedPassword); err != nil {
		logger.Info("User provided password does not match the hash stored in the database", "err", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expiresIn := time.Second * time.Duration(tokenExpirySeconds)
	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, expiresIn)
	if err != nil {
		logger.Info("Error issuing user token", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Error issuing user token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
