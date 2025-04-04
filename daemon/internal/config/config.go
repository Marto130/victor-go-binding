package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	Port        string
	DatabaseURL string
}

func LoadConfig() (*Config, error) {

	currentDir, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}
	envPaths := []string{
		filepath.Join(currentDir, ".env"),
	}

	var loaded bool
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			fmt.Printf("Config updated from %s", path)
			loaded = true
			break
		}
	}

	if configPath := os.Getenv("VICTOR_CONFIG_PATH"); configPath != "" {
		if err := godotenv.Load(configPath); err == nil {
			fmt.Printf("Configuration updated from %s (VICTOR_CONFIG_PATH)", configPath)
			loaded = true
		}
	}

	if !loaded {
		fmt.Println(".env file not found, using default values")
	}

	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "localhost")
	databaseURL := getEnv("DATABASE_URL", fmt.Sprintf("localhost:%s", port))

	return &Config{
		Host:        host,
		Port:        port,
		DatabaseURL: databaseURL,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
