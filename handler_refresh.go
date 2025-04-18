package main

import (
	"fmt"
	"net/http"

	"github.com/szmktk/chirpy/internal/auth"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		logger.Info("Error reading refresh token", "err", err)
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error reading refresh token: %s", err))
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret, tokenExpirationTime)
	if err != nil {
		logger.Info("Error issuing user token", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Error issuing user token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cfg.db.UpdateRefreshToken(r.Context(), token)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusNoContent)
}
