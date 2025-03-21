package main

import (
	"fmt"
	"log"

	"github.com/Tech-Arch1tect/config"
)

// example config struct with validation
type AppConfig struct {
	AppName string `env:"APP_NAME" validate:"required,min=10"`
	Port    int    `env:"PORT" validate:"required,max=65535"`
}

// default values if not set in environment
func (c *AppConfig) SetDefaults() {
	c.AppName = "Default App"
	c.Port = 8080
}

func main() {
	var cfg AppConfig
	if err := config.Load(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Loaded configuration: %+v\n", cfg)
}
