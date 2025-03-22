package config

import (
	"os"
	"testing"
)

type InValidatorConfig struct {
	InField string `env:"TEST_IN" validate:"required,in=a|b|c"`
}

func (i *InValidatorConfig) SetDefaults() {
}

func TestInValidatorValid(t *testing.T) {
	os.Setenv("TEST_IN", "b") // "b" is in the allowed list: a, b, c.
	defer os.Unsetenv("TEST_IN")

	var cfg InValidatorConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("expected valid InField, got error: %v", err)
	}
}

func TestInValidatorInvalid(t *testing.T) {
	os.Setenv("TEST_IN", "d") // "d" is not in the allowed list.
	defer os.Unsetenv("TEST_IN")

	var cfg InValidatorConfig
	err := Load(&cfg)
	if err == nil {
		t.Fatal("expected error for invalid InField, got nil")
	}
	expectedErr := "validation error: field 'InField' must be one of the following values: a|b|c"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
