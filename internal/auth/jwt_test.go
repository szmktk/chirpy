package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/szmktk/chirpy/internal/auth"
)

func TestMakeJWT(t *testing.T) {
	tokenSecret := "TestSecret"
	userID := uuid.New()
	expiresIn := time.Hour

	// Generate JWT
	token, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err, "error should not occur when generating a valid JWT")
	assert.NotEmpty(t, token, "token should not be empty")

	// Parse the generated token to verify claims
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	assert.NoError(t, err, "error should not occur when parsing the generated token")
	assert.NotNil(t, parsedToken, "parsed token should not be nil")
	assert.True(t, parsedToken.Valid, "parsed token should be valid")

	// Verify claims
	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	assert.True(t, ok, "claims should be of type *jwt.RegisteredClaims")
	assert.Equal(t, "chirpy", claims.Issuer, "issuer should match")
	assert.Equal(t, userID.String(), claims.Subject, "subject should match userID")
}

func TestValidateJWT(t *testing.T) {
	tokenSecret := "TestSecret"
	userID := uuid.New()
	expiresIn := time.Hour

	validToken, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err, "error should not occur when generating a valid JWT")

	scenarios := []struct {
		name           string
		token          string
		tokenSecret    string
		expectedUserID uuid.UUID
		expectedError  error
	}{
		{
			name:           "ok",
			token:          validToken,
			tokenSecret:    "TestSecret",
			expectedUserID: userID,
			expectedError:  nil,
		},
		{
			name:           "invalid token",
			token:          "InvalidToken",
			tokenSecret:    "TestSecret",
			expectedUserID: uuid.Nil,
			expectedError:  jwt.ErrTokenMalformed,
		},
		{
			name:           "wrong secret",
			token:          validToken,
			tokenSecret:    "WrongSecret",
			expectedUserID: uuid.Nil,
			expectedError:  jwt.ErrTokenSignatureInvalid,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			parsedUserID, err := auth.ValidateJWT(scenario.token, scenario.tokenSecret)

			if scenario.expectedError == nil {
				assert.NoError(t, err, "error should not occur when validating a valid JWT")
				assert.Equal(t, userID, parsedUserID, "parsed userID should match the original userID")
			} else {
				assert.Error(t, err, "error should occur when validating a token")
				assert.ErrorIs(t, err, scenario.expectedError)
			}
		})
	}
}

func TestMakeJWT_Expiration(t *testing.T) {
	tokenSecret := "TestSecret"
	userID := uuid.New()
	// Use short expiration time
	expiresIn := time.Millisecond * 10

	// Generate JWT
	token, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err, "error should not occur when generating a valid JWT")

	// Wait for the token to expire
	time.Sleep(time.Millisecond * 20)

	// Validate the expired JWT
	_, err = auth.ValidateJWT(token, tokenSecret)
	assert.Error(t, err, "error should occur when validating an expired JWT")
}

func TestGetBearerToken(t *testing.T) {
	scenarios := []struct {
		name          string
		headers       http.Header
		expectedToken string
		expectedError string
	}{
		{
			name:          "ok",
			headers:       http.Header{"Authorization": {"Bearer myJsonWebToken"}},
			expectedToken: "myJsonWebToken",
		},
		{
			name:          "missing auth header",
			headers:       http.Header{},
			expectedToken: "",
			expectedError: "authorization header not found",
		},
		{
			name:          "empty auth header",
			headers:       http.Header{"Authorization": {""}},
			expectedToken: "",
			expectedError: "authorization header not found",
		},
		{
			name:          "missing bearer prefix",
			headers:       http.Header{"Authorization": {"myJsonWebToken"}},
			expectedToken: "",
			expectedError: "invalid authorization header format",
		},
		{
			name:          "invalid auth header format",
			headers:       http.Header{"Authorization": {"BearermyJsonWebToken"}},
			expectedToken: "",
			expectedError: "invalid authorization header format",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			token, err := auth.GetBearerToken(scenario.headers)

			if scenario.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, scenario.expectedToken, token)
			} else {
				assert.Error(t, err, "error should occur when getting a bearer token")
				assert.EqualError(t, err, scenario.expectedError)
			}
		})
	}
}
