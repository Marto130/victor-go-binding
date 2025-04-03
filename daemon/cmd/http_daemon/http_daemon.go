package http_daemon

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"victorgo/daemon/internal/api"
)

type Server struct {
	httpServer *http.Server
	port       int
}

func NewServer(port int) *Server {

	router := api.SetupRouter()

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		},
		port: port,
	}
}

func (s *Server) Start() error {
	log.Printf("Initializing server on PORT: %d", s.port)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error server initialization: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
