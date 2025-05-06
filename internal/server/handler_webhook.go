package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/auth"
)

func (srv *Server) UpgradeUserWebhook(w http.ResponseWriter, r *http.Request) error {
	type input struct {
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
		Event string `json:"event"`
	}

	if apiKey, err := auth.GetApiKey(r.Header); err != nil || apiKey != srv.cfg.PolkaKey {
		srv.logger.Debug("Error getting api key", "err", err)
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	payload := input{}
	err := decoder.Decode(&payload)
	if err != nil {
		srv.logger.Error("Error decoding JSON body", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}
	srv.logger.Info("decoded payload", "json", payload)

	if payload.Event != "user.upgraded" {
		respondWithNoContent(w)
		return nil
	}

	srv.logger.Info("handling webhook event", ".event", payload.Event)
	_, err = srv.db.UpgradeUser(r.Context(), payload.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			srv.logger.Info("User not found", "user_id", payload.Data.UserID)
			return APIError{Status: http.StatusNotFound, Msg: "User with given id has not been found"}
		}
		srv.logger.Error("Error upgrading user", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	respondWithNoContent(w)
	return nil
}
