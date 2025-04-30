package server

import (
	"context"
	"net/http"

	"github.com/szmktk/chirpy/internal/auth"
)

type contextKey string

const contextKeyUserID contextKey = "userID"

// authMiddleware extracts and validates the Bearer JWT token, storing the user ID in context.
func (srv *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		userID, err := auth.ValidateJWT(token, srv.cfg.TokenSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
		next(w, r.WithContext(ctx))
	}
}
