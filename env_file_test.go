package config

import (
	"os"
	"testing"
)

type EnvFileConfig struct {
	AppName     string `env:"APP_NAME" validate:"required"`
	Port        int    `env:"PORT" validate:"min=1000"`
	Debug       bool   `env:"DEBUG"`
	DatabaseURL string `env:"DATABASE_URL"`
}

func (c *EnvFileConfig) SetDefaults() {
	c.AppName = "DefaultApp"
	c.Port = 8080
	c.Debug = false
	c.DatabaseURL = "localhost:5432"
}

func TestLoadFromEnvFile(t *testing.T) {
	envVars := []string{"APP_NAME", "PORT", "DEBUG", "DATABASE_URL"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	envContent := `# This is a comment
APP_NAME=TestApp
PORT=3000
DEBUG=true
DATABASE_URL="postgres://user:pass@localhost/db"

# Another comment
`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	var cfg EnvFileConfig
	err = Load(&cfg)
	if err != nil {
		t.Fatalf("Expected successful load from .env file, got error: %v", err)
	}

	if cfg.AppName != "TestApp" {
		t.Errorf("Expected AppName to be 'TestApp', got '%s'", cfg.AppName)
	}
	if cfg.Port != 3000 {
		t.Errorf("Expected Port to be 3000, got %d", cfg.Port)
	}
	if cfg.Debug != true {
		t.Errorf("Expected Debug to be true, got %v", cfg.Debug)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost/db" {
		t.Errorf("Expected DatabaseURL to be 'postgres://user:pass@localhost/db', got '%s'", cfg.DatabaseURL)
	}
}

func TestEnvironmentOverridesEnvFile(t *testing.T) {
	envVars := []string{"APP_NAME", "PORT", "DEBUG", "DATABASE_URL"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	envContent := `APP_NAME=EnvFileApp
PORT=3000
DEBUG=false
DATABASE_URL=env_file_db`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	os.Setenv("APP_NAME", "EnvironmentApp")
	os.Setenv("DEBUG", "true")
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("DEBUG")
	}()

	var cfg EnvFileConfig
	err = Load(&cfg)
	if err != nil {
		t.Fatalf("Expected successful load, got error: %v", err)
	}

	if cfg.AppName != "EnvironmentApp" {
		t.Errorf("Expected AppName to be 'EnvironmentApp' (from env), got '%s'", cfg.AppName)
	}
	if cfg.Debug != true {
		t.Errorf("Expected Debug to be true (from env), got %v", cfg.Debug)
	}

	if cfg.Port != 3000 {
		t.Errorf("Expected Port to be 3000 (from .env file), got %d", cfg.Port)
	}
	if cfg.DatabaseURL != "env_file_db" {
		t.Errorf("Expected DatabaseURL to be 'env_file_db' (from .env file), got '%s'", cfg.DatabaseURL)
	}
}

func TestMissingEnvFileDoesNotFail(t *testing.T) {
	os.Remove(".env")

	envVars := []string{"APP_NAME", "PORT", "DEBUG", "DATABASE_URL"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	var cfg EnvFileConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("Expected successful load even without .env file, got error: %v", err)
	}

	if cfg.AppName != "DefaultApp" {
		t.Errorf("Expected AppName to be 'DefaultApp' (default), got '%s'", cfg.AppName)
	}
	if cfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080 (default), got %d", cfg.Port)
	}
}

func TestEnvFileWithVariousFormats(t *testing.T) {
	envVars := []string{"VAR1", "VAR2", "VAR3", "VAR4", "VAR5"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	envContent := `# Comment at the beginning
VAR1=simple_value
VAR2="quoted value"
VAR3='single quoted'
VAR4=value with spaces
VAR5=

# Comment in the middle

# Empty line above
VAR6=another_value
MALFORMED_LINE_NO_EQUALS
=VALUE_WITHOUT_KEY
`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	err = loadEnvFile(".env")
	if err != nil {
		t.Fatalf("Expected successful .env load, got error: %v", err)
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"VAR1", "simple_value"},
		{"VAR2", "quoted value"},
		{"VAR3", "single quoted"},
		{"VAR4", "value with spaces"},
		{"VAR6", "another_value"},
	}

	for _, test := range tests {
		if value, exists := os.LookupEnv(test.key); !exists {
			t.Errorf("Expected environment variable %s to be set", test.key)
		} else if value != test.expected {
			t.Errorf("Expected %s to be '%s', got '%s'", test.key, test.expected, value)
		}
	}

	if value, exists := os.LookupEnv("VAR5"); !exists {
		t.Errorf("Expected VAR5 to be set (even if empty)")
	} else if value != "" {
		t.Errorf("Expected VAR5 to be empty, got '%s'", value)
	}
}
