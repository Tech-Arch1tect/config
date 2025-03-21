package config

import (
	"os"
	"testing"
)

type TestConfig struct {
	Name string `env:"TEST_NAME" validate:"required,min=5,max=20"`
	Port int    `env:"TEST_PORT" validate:"min=1024,max=65535"`
}

func (c *TestConfig) SetDefaults() {
	c.Name = "DefaultName"
	c.Port = 8080
}

func TestLoadValidConfig(t *testing.T) {
	os.Unsetenv("TEST_NAME")
	os.Unsetenv("TEST_PORT")

	os.Setenv("TEST_NAME", "ValidApp")
	os.Setenv("TEST_PORT", "9000")

	var cfg TestConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}

	if cfg.Name != "ValidApp" {
		t.Errorf("expected Name to be 'ValidApp', got '%s'", cfg.Name)
	}
	if cfg.Port != 9000 {
		t.Errorf("expected Port to be 9000, got %d", cfg.Port)
	}
}

// TestLoadMissingRequired tests that missing required values trigger a validation error.
func TestLoadMissingRequired(t *testing.T) {
	os.Unsetenv("TEST_NAME")
	os.Unsetenv("TEST_PORT")

	var cfg TestConfig

	cfg.SetDefaults()
	cfg.Name = ""

	err := ValidateStruct(&cfg)
	if err == nil {
		t.Fatal("expected error for missing required field Name, got nil")
	}

	expectedErr := "field 'Name' is required"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

// TestLoadValidationRules tests the min and max validations for both string and integer fields.
func TestLoadValidationRules(t *testing.T) {
	tests := []struct {
		desc      string
		name      string
		port      int
		expectErr bool
		errMsg    string
	}{
		{
			desc:      "Name too short",
			name:      "App", // less than min length 5
			port:      9000,
			expectErr: true,
			errMsg:    "field 'Name' must be at least 5 characters",
		},
		{
			desc:      "Port too low",
			name:      "ValidApp",
			port:      800, // below minimum 1024
			expectErr: true,
			errMsg:    "field 'Port' must be at least 1024",
		},
		{
			desc:      "Port too high",
			name:      "ValidApp",
			port:      70000, // above maximum 65535
			expectErr: true,
			errMsg:    "field 'Port' must be at most 65535",
		},
		{
			desc:      "Valid config using defaults",
			name:      "ValidConfig",
			port:      8080,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cfg := TestConfig{
				Name: tt.name,
				Port: tt.port,
			}
			err := ValidateStruct(&cfg)
			if tt.expectErr {
				if err == nil {
					t.Fatal("expected an error but got nil")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, but got: %v", err)
				}
			}
		})
	}
}
