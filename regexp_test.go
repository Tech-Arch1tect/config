package config

import (
	"os"
	"testing"
)

type RegexConfig struct {
	Username string `env:"TEST_USERNAME" validate:"required,regexp=^[a-zA-Z0-9]+$"`
}

func (r *RegexConfig) SetDefaults() {
}

func TestRegexpValidatorValid(t *testing.T) {
	os.Setenv("TEST_USERNAME", "User123")
	defer os.Unsetenv("TEST_USERNAME")

	var cfg RegexConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("expected valid username, got error: %v", err)
	}
}

func TestRegexpValidatorInvalid(t *testing.T) {
	os.Setenv("TEST_USERNAME", "User 123") // Contains a space, invalid per regex.
	defer os.Unsetenv("TEST_USERNAME")

	var cfg RegexConfig
	err := Load(&cfg)
	if err == nil {
		t.Fatal("expected error for invalid username, got nil")
	}
	expectedErr := "validation error: field 'Username' must match the pattern '^[a-zA-Z0-9]+$'"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
