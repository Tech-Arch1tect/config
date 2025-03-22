package config

import (
	"os"
	"testing"
)

type URLConfig struct {
	URL string `env:"TEST_URL" validate:"required,url"`
}

func (u *URLConfig) SetDefaults() {
}

func TestURLValidatorValid(t *testing.T) {
	os.Setenv("TEST_URL", "https://example.com")
	defer os.Unsetenv("TEST_URL")

	var cfg URLConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("expected valid URL, got error: %v", err)
	}
}

func TestURLValidatorInvalid(t *testing.T) {
	os.Setenv("TEST_URL", "not a url")
	defer os.Unsetenv("TEST_URL")

	var cfg URLConfig
	err := Load(&cfg)
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
	expectedErr := "validation error: field 'URL' must be a valid URL"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
