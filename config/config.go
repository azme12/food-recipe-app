package config

import (
    "os"
    "fmt"
)

// Config represents the application configuration structure
type Config struct {
    DatabaseURL string `json:"DatabaseURL"`
    Port        string `json:"Port"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    dbPassword := os.Getenv("DB_PASSWORD")
    if dbPassword == "" {
        return nil, fmt.Errorf("DB_PASSWORD environment variable not set")
    }

    // Load other configuration parameters
    config := &Config{
        DatabaseURL: "postgres://postgres:" + dbPassword + "@localhost:5432/food_recipes",
        Port:        os.Getenv("PORT"),
    }

    return config, nil
}
