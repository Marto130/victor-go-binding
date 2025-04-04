package main

import (
	"log"
	"net/http"
	routes "victorgo/daemon/internal/api"
	config "victorgo/daemon/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize the HTTP server
	router := routes.SetupRouter()

	// Start the server
	address := cfg.Host + ":" + cfg.Port
	log.Printf("Victor daemon running on Port: %s", cfg.Port)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("Error starting Victor daemon: %v", err)
	}
}
