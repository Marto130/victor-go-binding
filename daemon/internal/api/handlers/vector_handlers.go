package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"victorgo/daemon/pkg/models"
	"victorgo/daemon/pkg/store"

	"github.com/gorilla/mux"
)

// InsertVector handles the insertion of a vector into the index
func InsertVectorHandler(w http.ResponseWriter, r *http.Request) {
	var req models.InsertVectorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Error decoding request:", err)
		http.Error(w, "Invalid insert vector request payload", http.StatusBadRequest)
		return
	}

	indexIDParam := r.URL.Query().Get("index_id")

	indexResource, exists := store.GetIndex(indexIDParam)
	if !exists {
		http.Error(w, "Index not found", http.StatusNotFound)
		return
	}

	vIndex := indexResource.VIndex

	dims, dimsExists := store.GetIndexDims(indexIDParam)
	if !dimsExists {
		http.Error(w, "Index dimensions not found", http.StatusNotFound)
		return
	}
	if len(req.Vector) != int(dims) {
		http.Error(w, "Vector dimensions do not match index dimensions", http.StatusBadRequest)
		return
	}

	if err := vIndex.Insert(req.ID, req.Vector); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{
		Status:  "success",
		Message: "Vector inserted successfully",
	})

}

// SearchVector handles searching for a vector in the index
func SearchVectorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	indexID := vars["indexID"]

	indexResource, exists := store.GetIndex(indexID)
	if !exists {
		http.Error(w, "Index not found", http.StatusNotFound)
		return
	}

	vIndex := indexResource.VIndex

	vectorParam := r.URL.Query().Get("vector")
	if vectorParam == "" {
		http.Error(w, "Missing vector parameter", http.StatusBadRequest)
		return
	}

	vectorStrings := strings.Split(vectorParam, ",")
	vector := make([]float32, len(vectorStrings))
	for i, s := range vectorStrings {
		val, err := strconv.ParseFloat(s, 32)
		if err != nil {
			http.Error(w, "Invalid vector value: "+s, http.StatusBadRequest)
			return
		}
		vector[i] = float32(val)
	}

	result, err := vIndex.Search(vector, len(vector))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{
		Status: "success",
		Data:   result,
	})
}
