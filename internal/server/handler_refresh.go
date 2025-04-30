package server

import (
	"net/http"

	"time"

	"github.com/szmktk/chirpy/internal/auth"
)

func (srv *Server) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	refreshToken, err := srv.db.GetRefreshToken(r.Context(), token)
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

	accessToken, err := auth.MakeJWT(refreshToken.UserID, srv.cfg.TokenSecret, accessTokenExpirationTime)
	if err != nil {
		srv.logger.Error("Error issuing user token", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (srv *Server) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := srv.db.RevokeRefreshToken(r.Context(), token); err != nil {
		srv.logger.Error("Error revoking refresh token", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithNoContent(w)
}
