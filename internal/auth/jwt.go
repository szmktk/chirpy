package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const tokenIssuer string = "chirpy"

// MakeJWT generates a signed JWT for a given user ID using the provided secret key.
// Returns the signed token string or an error if signing fails.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	return token.SignedString(signingKey)
}

// ValidateJWT parses and validates a JWT string using the provided secret.
// Returns the parsed user ID or an error if validation fails.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != tokenIssuer {
		return uuid.Nil, errors.New("invalid issuer")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return userUUID, nil
}

// extractAuthToken extracts a token for the given scheme from the HTTP Authorization header.
// Returns the token string or an error if the header is missing or malformed.
func extractAuthToken(headers http.Header, scheme string) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != scheme {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}

// GetBearerToken extracts the bearer token from an HTTP Authorization header.
// Returns the token string or an error if the header is missing or malformed.
func GetBearerToken(headers http.Header) (string, error) {
	return extractAuthToken(headers, "bearer")
}

// GetApiKey extracts the API key from an HTTP Authorization header.
// Returns the key string or an error if the header is missing or malformed.
func GetApiKey(headers http.Header) (string, error) {
	return extractAuthToken(headers, "apikey")
}

// MakeRefreshToken generates a new secure random refresh token.
// Returns the token or an error if random data generation fails.
func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}
