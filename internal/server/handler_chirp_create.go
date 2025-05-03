package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/database"
)

const maxChirpLength int = 140

var forbiddenWords = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (srv *Server) CreateChirp(w http.ResponseWriter, r *http.Request) error {
	type input struct {
		Body string `json:"body"`
	}

	ctxVal := r.Context().Value(contextKeyUserID)
	parsedUserID, ok := ctxVal.(uuid.UUID)
	if !ok {
		return APIError{Status: http.StatusUnauthorized, Msg: "Unauthorized"}
	}

	decoder := json.NewDecoder(r.Body)
	payload := input{}
	err := decoder.Decode(&payload)
	if err != nil {
		srv.logger.Error("Error decoding JSON body: %s", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	if len(payload.Body) > maxChirpLength {
		return APIError{Status: http.StatusBadRequest, Msg: "Chirp is too long"}
	}

	cleanedBody := sanitizeInput(payload.Body)

	chirp, err := srv.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: parsedUserID,
	})
	if err != nil {
		srv.logger.Error("Error creating chirp: %s", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}

	return respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func sanitizeInput(body string) string {
	var sanitized []string
	for _, word := range strings.Split(body, " ") {
		if forbiddenWords[strings.ToLower(word)] {
			word = "****"
		}
		sanitized = append(sanitized, word)
	}
	return strings.Join(sanitized, " ")
}
