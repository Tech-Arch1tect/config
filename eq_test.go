package config

import (
	"os"
	"testing"
)

type EqValidatorConfig struct {
	EqField string `env:"TEST_EQ" validate:"required,eq=hello"`
}

func (e *EqValidatorConfig) SetDefaults() {
}

func TestEqValidatorValid(t *testing.T) {
	os.Setenv("TEST_EQ", "hello")
	defer os.Unsetenv("TEST_EQ")

	var cfg EqValidatorConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("expected valid EqField, got error: %v", err)
	}
}

func TestEqValidatorInvalid(t *testing.T) {
	os.Setenv("TEST_EQ", "world")
	defer os.Unsetenv("TEST_EQ")

	var cfg EqValidatorConfig
	err := Load(&cfg)
	if err == nil {
		t.Fatal("expected error for invalid EqField, got nil")
	}
	expectedErr := "validation error: field 'EqField' must be equal to 'hello'"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
