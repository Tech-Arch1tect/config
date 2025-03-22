package config

import (
	"os"
	"testing"
)

type NeValidatorConfig struct {
	NeField string `env:"TEST_NE" validate:"required,ne=forbidden"`
}

func (n *NeValidatorConfig) SetDefaults() {
}

func TestNeValidatorValid(t *testing.T) {
	os.Setenv("TEST_NE", "allowed") // "allowed" is not equal to "forbidden"
	defer os.Unsetenv("TEST_NE")

	var cfg NeValidatorConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("expected valid NeField, got error: %v", err)
	}
}

func TestNeValidatorInvalid(t *testing.T) {
	os.Setenv("TEST_NE", "forbidden") // "forbidden" is equal to the disallowed value
	defer os.Unsetenv("TEST_NE")

	var cfg NeValidatorConfig
	err := Load(&cfg)
	if err == nil {
		t.Fatal("expected error for invalid NeField, got nil")
	}
	expectedErr := "validation error: field 'NeField' must not be equal to 'forbidden'"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
