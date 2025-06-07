# github.com/Tech-Arch1tect/config

An extremely simple library to load configurations from environment variables and .env files.

## Features

- **Environment Variable Loading**: Load configuration from environment variables
- **.env File Support**: Automatically loads variables from a .env file
- **Environment Override**: Environment variables take precedence over .env file values
- **Validation**: ~~Comprehensive~~ (not yet!) validation rules for configuration values
- **Default Values**: Set default values through the `SetDefaults` method
- **No External Dependencies**: Pure Go implementation

## Supported Validations (more to come)

- **required**: Ensures a value is provided.
- **min**: Checks that a numeric value is greater than or equal to a minimum, or that a string has at least the minimum number of characters.
- **max**: Checks that a numeric value is less than or equal to a maximum, or that a string does not exceed the maximum number of characters.
- **email**: Validates that a string is a properly formatted email address.
- **url**: Validates that a string is a properly formatted URL.
- **regexp**: Validates that a string matches a given regular expression pattern.
- **in**: Validates that a string is one of a set of allowed values, using the pipe (`|`) character as a delimiter.
- **not_in**: Validates that a string is not one of a set of disallowed values, using the pipe (`|`) character as a delimiter.
- **eq**: Validates that a string is equal to a specified value.
- **ne**: Validates that a string is not equal to a specified value.

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
    Port    int    `env:"PORT" validate:"required,max=65535"`
    Email   string `env:"APP_EMAIL" validate:"required,email"`
}

func (c *AppConfig) SetDefaults() {
    c.AppName = "Default App"
    c.Port = 8080
    c.Email = "default@example.com"
}

func main() {
    var cfg AppConfig
    if err := config.Load(&cfg); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    fmt.Printf("Loaded configuration: %+v\n", cfg)
}
```

## .env File Support

The library automatically loads environment variables from a `.env` file in the current directory. The .env file format supports:

- Simple key=value pairs
- Quoted values (both single and double quotes)
- Comments (lines starting with #)
- Empty lines (ignored)

### Example .env file:

```bash
# Application Configuration
APP_NAME="My Application"
PORT=3000
DEBUG=true
DATABASE_URL='postgres://user:pass@localhost/mydb'

# Email Configuration
APP_EMAIL=admin@example.com
```

### Priority Order

Configuration values are loaded in the following priority order (highest to lowest):

1. **Environment Variables** - Values set in the actual environment
2. **.env File** - Values from the .env file
3. **Default Values** - Values set in the `SetDefaults()` method

This means environment variables will always override .env file values, and .env file values will override defaults.

### Example with .env file:

**`.env` file:**

```bash
APP_NAME=EnvFileApp
PORT=3000
```

**Environment variable:**

```bash
export APP_NAME=EnvironmentApp
```

**Result:**

- `APP_NAME` will be `"EnvironmentApp"` (from environment variable)
- `PORT` will be `3000` (from .env file)
- `EMAIL` will be `"default@example.com"` (from defaults)

## How It Works

- **.env File Loading:**  
  The library first reads a `.env` file from the current directory (if it exists). Variables from this file are set in the environment only if they don't already exist.

- **Environment Variable Loading:**  
  The library scans your struct for `env` tags and assigns the corresponding environment variable values. If an environment variable is not set, the `SetDefaults` method provides fallback values.

- **Validation:**  
  After loading the configuration, the library validates the struct using the `validate` tags described above.

## Error Handling

- If the `.env` file doesn't exist, the library continues without error
- If the `.env` file exists but has parsing errors, the library will return an error
- Validation errors are returned if any configured validation rules fail
