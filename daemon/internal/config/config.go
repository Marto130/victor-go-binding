package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds the configuration settings for the daemon
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
			log.Printf("Configuración cargada desde %s", path)
			loaded = true
			break
		}
	}

	if configPath := os.Getenv("VICTOR_CONFIG_PATH"); configPath != "" {
		if err := godotenv.Load(configPath); err == nil {
			log.Printf("Configuración cargada desde %s (VICTOR_CONFIG_PATH)", configPath)
			loaded = true
		}
	}

	if !loaded {
		log.Println("No se encontró archivo .env, usando valores predeterminados")
	}

	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "localhost")
	databaseURL := getEnv("DATABASE_URL", "localhost:5432")

	return &Config{
		Host:        host,
		Port:        port,
		DatabaseURL: databaseURL,
	}, nil
}

// getEnv retrieves the value of the environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
