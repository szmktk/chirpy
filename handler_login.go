package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/szmktk/chirpy/internal/auth"
	"github.com/szmktk/chirpy/internal/database"
)

const accessTokenExpirationTime time.Duration = time.Hour

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := input{}
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

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, accessTokenExpirationTime)
	if err != nil {
		logger.Info("Error issuing user token", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Error issuing user token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		logger.Info("Error issuing refresh token", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Error issuing refresh token")
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error saving refresh token: %s", err))
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
		Token:        token,
		RefreshToken: refreshToken,
	})
}
