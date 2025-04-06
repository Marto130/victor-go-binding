package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"victorgo/daemon/internal/api"
	"victorgo/daemon/internal/config"
)

func main() {

	flag.Parse()

	router := *api.SetupRouter()

	address := *config.Host + ":" + strconv.Itoa(*config.Port)
	log.Printf("Victor daemon running on Host: %s Port: %s", *config.Host, strconv.Itoa(*config.Port))
	if err := http.ListenAndServe(address, &router); err != nil {
		log.Fatalf("Error starting Victor daemon: %v", err)
	}
}
