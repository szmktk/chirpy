package main

import (
	"net/http"

	"time"

	"github.com/szmktk/chirpy/internal/auth"
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
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret, accessTokenExpirationTime)
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

	if err := cfg.db.RevokeRefreshToken(r.Context(), token); err != nil {
		logger.Info("Error revoking refresh token", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token")
		return
	}

	respondWithNoContent(w)
}
