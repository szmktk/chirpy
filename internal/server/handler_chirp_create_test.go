package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/szmktk/chirpy/internal/config"
	"github.com/szmktk/chirpy/internal/database"
)

func TestCreateChirp(t *testing.T) {
	tests := []struct {
		name           string
		userID         uuid.UUID
		requestBody    string
		expectedStatus int
		expectedError  string
		wantBody       string
	}{
		{
			name:           "happy path",
			userID:         uuid.New(),
			requestBody:    `{"body": "Hello world!"}`,
			expectedStatus: http.StatusCreated,
			wantBody:       "Hello world!",
		},
		{
			name:           "unauthorized - no user ID",
			requestBody:    `{"body": "Hello world!"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
		{
			name:           "invalid JSON body",
			userID:         uuid.New(),
			requestBody:    `invalid json`,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal Server Error",
		},
		{
			name:           "empty JSON body",
			userID:         uuid.New(),
			requestBody:    "",
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal Server Error",
		},
		{
			name:           "chirp too long",
			userID:         uuid.New(),
			requestBody:    `{"body": "` + strings.Repeat("a", maxChirpLength+1) + `"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Chirp is too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cfg := &config.Config{}
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			defer db.Close()

			queries := database.New(db)
			srv, err := NewServer(cfg, queries, logger)
			if err != nil {
				t.Fatalf("failed to create server: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/chirps", strings.NewReader(tt.requestBody))

			// Only set up mock expectations for the happy path
			if tt.userID != uuid.Nil && tt.expectedError == "" {
				ctx := context.WithValue(req.Context(), contextKeyUserID, tt.userID)
				req = req.WithContext(ctx)

				// Setup mock DB expectations for successful case
				mock.ExpectQuery("INSERT INTO chirps").
					WithArgs(tt.wantBody, tt.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "body", "user_id"}).
						AddRow(uuid.New(), time.Now(), time.Now(), tt.wantBody, tt.userID))
			} else if tt.userID != uuid.Nil {
				ctx := context.WithValue(req.Context(), contextKeyUserID, tt.userID)
				req = req.WithContext(ctx)
			}

			// Execute request
			w := httptest.NewRecorder()
			err = srv.CreateChirp(w, req)

			// Verify response
			if tt.expectedError != "" {
				apiErr, ok := err.(APIError)
				if !ok {
					t.Errorf("expected APIError, got %T", err)
					return
				}
				if apiErr.Status != tt.expectedStatus {
					t.Errorf("expected status %d, got %d", tt.expectedStatus, apiErr.Status)
				}
				if apiErr.Msg != tt.expectedError {
					t.Errorf("expected error message %q, got %q", tt.expectedError, apiErr.Msg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if w.Code != tt.expectedStatus {
					t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
				}

				var response Chirp
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
					return
				}
				if response.Body != tt.wantBody {
					t.Errorf("expected body %q, got %q", tt.wantBody, response.Body)
				}
				if response.UserID != tt.userID {
					t.Errorf("expected user ID %v, got %v", tt.userID, response.UserID)
				}
			}

			// Only verify mock expectations for the happy path
			if tt.expectedError == "" {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("unmet mock expectations: %v", err)
				}
			}
		})
	}
}
