package server

import (
	"net/http"

	"time"

	"github.com/szmktk/chirpy/internal/auth"
)

func (srv *Server) Refresh(w http.ResponseWriter, r *http.Request) error {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	refreshToken, err := srv.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if refreshToken.RevokedAt.Valid {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, srv.cfg.TokenSecret, accessTokenExpirationTime)
	if err != nil {
		srv.logger.Error("Error issuing user token", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	return respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (srv *Server) Revoke(w http.ResponseWriter, r *http.Request) error {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	if err := srv.db.RevokeRefreshToken(r.Context(), token); err != nil {
		srv.logger.Error("Error revoking refresh token", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	respondWithNoContent(w)
	return nil
}
