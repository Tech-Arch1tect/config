# github.com/Tech-Arch1tect/config

An extremely simple library to load configurations from environment variables.

## Supported Validations (more to come)

- **required**: Ensures a value is provided.
- **min**: Checks that a numeric value is greater than or equal to a minimum, or that a string has at least the minimum number of characters.
- **max**: Checks that a numeric value is less than or equal to a maximum, or that a string does not exceed the maximum number of characters.

## Installation

```bash
go get github.com/Tech-Arch1tect/config
```

## Usage

Define your configuration struct and provide default values by implementing the `SetDefaults` method.

```go
package main

import (
    "fmt"
    "log"

    "github.com/Tech-Arch1tect/config"

)

type AppConfig struct {
AppName string `env:"APP_NAME" validate:"required,min=10"`
Port int `env:"PORT" validate:"required,max=65535"`
}

func (c \*AppConfig) SetDefaults() {
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
```

## How It Works

- **Environment Variable Loading:**  
  The library scans your struct for `env` tags and assigns the corresponding environment variable values. If an environment variable is not set, the `SetDefaults` method provides fallback values.

- **Validation:**  
  After loading the configuration, the library validates the struct using the `validate` tags described above.
