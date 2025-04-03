package api

import (
	"net/http"
	"victorgo/daemon/internal/api/handlers"
	"victorgo/daemon/pkg/routes"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc(routes.CreateIndexPath, handlers.CreateIndexHandler()).Methods(http.MethodPost)
	r.HandleFunc(routes.InsertVectorPath, handlers.InsertVectorHandler).Methods(http.MethodPost)
	r.HandleFunc(routes.SearchVectorPath, handlers.SearchVectorHandler).Methods(http.MethodGet)
}

func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	RegisterRoutes(router)
	return router
}
