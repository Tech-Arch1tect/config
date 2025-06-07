package main

import (
	"fmt"
	"log"

	"github.com/Tech-Arch1tect/config"
)

// example config struct with validation
type AppConfig struct {
	AppName     string `env:"APP_NAME" validate:"required,min=10"`
	Port        int    `env:"PORT" validate:"required,max=65535"`
	Debug       bool   `env:"DEBUG"`
	DatabaseURL string `env:"DATABASE_URL" validate:"required,url"`
	APIKey      string `env:"API_KEY" validate:"required,min=10"`
	SMTPHost    string `env:"SMTP_HOST"`
	SMTPPort    int    `env:"SMTP_PORT" validate:"min=1,max=65535"`
}

// default values if not set in environment
func (c *AppConfig) SetDefaults() {
	c.AppName = "Default App Name"
	c.Port = 8080
	c.Debug = false
	c.DatabaseURL = "postgres://localhost:5432/defaultdb"
	c.APIKey = "default-api-key"
	c.SMTPHost = "localhost"
	c.SMTPPort = 587
}

func main() {
	fmt.Println("Loading configuration...")
	fmt.Println("Priority order: Environment Variables > .env file > Defaults")
	fmt.Println()

	var cfg AppConfig
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Loaded configuration:\n")
	fmt.Printf("  App Name:     %s\n", cfg.AppName)
	fmt.Printf("  Port:         %d\n", cfg.Port)
	fmt.Printf("  Debug:        %v\n", cfg.Debug)
	fmt.Printf("  Database URL: %s\n", cfg.DatabaseURL)
	fmt.Printf("  API Key:      %s\n", cfg.APIKey)
	fmt.Printf("  SMTP Host:    %s\n", cfg.SMTPHost)
	fmt.Printf("  SMTP Port:    %d\n", cfg.SMTPPort)

	fmt.Println()
	fmt.Println("Try the following to see the priority in action:")
	fmt.Println("1. Run as-is (uses .env file and defaults)")
	fmt.Println("2. Set an environment variable: export APP_NAME=\"Environment Override\"")
	fmt.Println("3. Run again to see the environment variable override the .env file")
}
