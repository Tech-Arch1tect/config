package config

import (
	"os"
	"testing"
)

type NotInValidatorConfig struct {
	NotInField string `env:"TEST_NOTIN" validate:"required,not_in=a|b|c"`
}

func (n *NotInValidatorConfig) SetDefaults() {
}

func TestNotInValidatorValid(t *testing.T) {
	// "d" is not in the disallowed list: a, b, c.
	os.Setenv("TEST_NOTIN", "d")
	defer os.Unsetenv("TEST_NOTIN")

	var cfg NotInValidatorConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("expected valid NotInField, got error: %v", err)
	}
}

func TestNotInValidatorInvalid(t *testing.T) {
	// "a" is in the disallowed list: a, b, c.
	os.Setenv("TEST_NOTIN", "a")
	defer os.Unsetenv("TEST_NOTIN")

	var cfg NotInValidatorConfig
	err := Load(&cfg)
	if err == nil {
		t.Fatal("expected error for invalid NotInField, got nil")
	}
	expectedErr := "validation error: field 'NotInField' must not be one of the following values: a|b|c"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
