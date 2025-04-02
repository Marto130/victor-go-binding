package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
	}
	s.initializeRoutes()
	return s
}

func (s *Server) initializeRoutes() {
	// Define your routes here
	// Example: s.router.HandleFunc("/api/index", indexHandler).Methods("POST")
}

func (s *Server) Start(addr string) {
	srv := &http.Server{
		Handler:      s.router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting server on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}