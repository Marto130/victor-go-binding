package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"victorgo/daemon/pkg/models"
	"victorgo/daemon/pkg/store"

	binding "victorgo/binding"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	indexMutex sync.RWMutex
)

func CreateIndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var createIndexRequest models.CreateIndexRequest

		urlParams := mux.Vars(r)
		indexNameParam := urlParams["indexName"]
		if indexNameParam == "" {
			http.Error(w, "Index name is required", http.StatusBadRequest)
			return
		}

		indexMutex.Lock()
		defer indexMutex.Unlock()

		if err := json.NewDecoder(r.Body).Decode(&createIndexRequest); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		index, err := binding.AllocIndex(createIndexRequest.IndexType, createIndexRequest.Method, createIndexRequest.Dims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		indexID := uuid.New().String()

		indexResource := models.IndexResource{
			CreateIndexRequest: models.CreateIndexRequest{
				IndexType: createIndexRequest.IndexType,
				Method:    createIndexRequest.Method,
				Dims:      createIndexRequest.Dims,
			},
			VIndex:    index,
			IndexName: indexNameParam,
			IndexID:   indexID,
		}

		store.StoreIndex(&indexResource)

		response := models.CreateIndexResponse{
			Status:  "success",
			Message: "Index created successfully",
			Results: models.CreateIndexResult{
				IndexName: indexNameParam,
				ID:        indexID,
				Dims:      createIndexRequest.Dims,
				IndexType: createIndexRequest.IndexType,
				Method:    createIndexRequest.Method,
			},
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
