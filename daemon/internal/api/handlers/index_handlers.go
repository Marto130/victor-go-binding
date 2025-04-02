package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	binding "victorgo/binding"
	"victorgo/daemon/pkg/models"
	"victorgo/daemon/pkg/store"

	"github.com/google/uuid"
)

var (
	indexMutex sync.RWMutex
)

// CreateIndexHandler handles the creation of a new index
func CreateIndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var createIndexRequest models.CreateIndexRequest

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
			IndexID: indexID,
			VIndex:  index,
		}

		store.StoreIndex(&indexResource)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Index created successfully", "id": indexID})
	}
}
