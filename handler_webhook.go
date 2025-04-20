package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUserWebhook(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type input struct {
		Data  `json:"data"`
		Event string `json:"event"`
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
	} else {
		logger.Info("handling webhook event", ".event", payload.Event)
		_, err := cfg.db.UpgradeUser(r.Context(), payload.UserID)
		if err != nil {
			logger.Info("User not found", "user_id", payload.UserID)
			respondWithError(w, http.StatusNotFound, "User with given id has not been found")
			return
		}

		respondWithNoContent(w)
	}
}
