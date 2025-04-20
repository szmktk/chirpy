package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUserWebhook(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
		Event string `json:"event"`
	}

	if apiKey, err := auth.GetApiKey(r.Header); err != nil || apiKey != cfg.polkaKey {
		logger.Debug("Error getting api key", "err", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	payload := input{}
	err := decoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding JSON body: %s", err))
		return
	}
	logger.Info("decoded payload", "json", payload)

	if payload.Event != "user.upgraded" {
		respondWithNoContent(w)
		return
	}

	logger.Info("handling webhook event", ".event", payload.Event)
	_, err = cfg.db.UpgradeUser(r.Context(), payload.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info("User not found", "user_id", payload.Data.UserID)
			respondWithError(w, http.StatusNotFound, "User with given id has not been found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upgrading user: %s", err))
		return
	}

	respondWithNoContent(w)
}
