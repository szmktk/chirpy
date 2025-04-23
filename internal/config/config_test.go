package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	os.Setenv("DB_URL", "postgres://localhost/db")
	os.Setenv("POLKA_KEY", "examplePolkaKey")
	os.Setenv("TOKEN_SECRET", "exampleTokenSecret")
	defer os.Unsetenv("DB_URL")
	defer os.Unsetenv("POLKA_KEY")
	defer os.Unsetenv("TOKEN_SECRET")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.DBURL != "postgres://localhost/db" {
		t.Errorf("expected DBURL to be 'postgres://localhost/db', got %s", config.DBURL)
	}

	if config.PolkaKey != "examplePolkaKey" {
		t.Errorf("expected PolkaKey to be 'examplePolkaKey', got %s", config.PolkaKey)
	}

	if config.TokenSecret != "exampleTokenSecret" {
		t.Errorf("expected TokenSecret to be 'exampleTokenSecret', got %s", config.TokenSecret)
	}
}

func TestLoadConfig_MissingVars(t *testing.T) {
	os.Unsetenv("DB_URL")
	os.Setenv("POLKA_KEY", "examplePolkaKey")
	os.Setenv("TOKEN_SECRET", "exampleTokenSecret")
	defer os.Unsetenv("POLKA_KEY")
	defer os.Unsetenv("TOKEN_SECRET")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected an error, got none")
	}

	expectedError := "missing required env vars: DB_URL"
	if err.Error() != expectedError {
		t.Errorf("expected error message to be '%s', got %s", expectedError, err.Error())
	}
}
