package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/szmktk/chirpy/internal/auth"
	"github.com/szmktk/chirpy/internal/database"
)

const accessTokenExpirationTime time.Duration = time.Hour

func (srv *Server) Login(w http.ResponseWriter, r *http.Request) error {
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
		srv.logger.Error("Error decoding JSON body: %s", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	if payload.Email == "" {
		return APIError{Status: http.StatusBadRequest, Msg: "Email cannot be empty"}
	}
	if payload.Password == "" {
		return APIError{Status: http.StatusBadRequest, Msg: "Password cannot be empty"}
	}

	user, err := srv.db.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		srv.logger.Info("User not found", "err", err)
		return APIError{Status: http.StatusUnauthorized, Msg: "Incorrect email or password"}
	}

	if err := auth.CheckPasswordHash(payload.Password, user.HashedPassword); err != nil {
		srv.logger.Info("User provided password does not match the hash stored in the database", "err", err)
		return APIError{Status: http.StatusUnauthorized, Msg: "Incorrect email or password"}
	}

	token, err := auth.MakeJWT(user.ID, srv.cfg.TokenSecret, accessTokenExpirationTime)
	if err != nil {
		srv.logger.Error("Error issuing user token", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		srv.logger.Error("Error issuing refresh token", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	_, err = srv.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		srv.logger.Error("Error saving refresh token: %s", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	return respondWithJSON(w, http.StatusOK, response{
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
