package api

import (
	"net/http"
	"victorgo/daemon/internal/api/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/index", handlers.CreateIndexHandler()).Methods(http.MethodPost)
	r.HandleFunc("/api/vector/{indexID}", handlers.InsertVectorHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/vector/{indexID}/search", handlers.SearchVectorHandler).Methods(http.MethodGet)
}

func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	RegisterRoutes(router)
	return router
}
