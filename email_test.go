package config

import (
	"os"
	"testing"
)

type EmailConfig struct {
	Email string `env:"TEST_EMAIL" validate:"required,email"`
}

func (e *EmailConfig) SetDefaults() {
}

func TestEmailValidatorValid(t *testing.T) {
	os.Setenv("TEST_EMAIL", "test@example.com")
	defer os.Unsetenv("TEST_EMAIL")

	var cfg EmailConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("expected valid email, got error: %v", err)
	}
}

func TestEmailValidatorInvalid(t *testing.T) {
	os.Setenv("TEST_EMAIL", "invalid-email")
	defer os.Unsetenv("TEST_EMAIL")

	var cfg EmailConfig
	err := Load(&cfg)
	if err == nil {
		t.Fatal("expected error for invalid email, got nil")
	}
	expectedErr := "validation error: field 'Email' must be a valid email address"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
