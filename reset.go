package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("Operation not allowed on platform: '%s'", cfg.platform))
		return
	}

	cfg.fileserverHits.Store(0)
	if err := cfg.db.DeleteUsers(r.Context()); err != nil {
		logger.Error("Error deleting users: %s", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hits reset to 0 and database reset to initial state"})
}
