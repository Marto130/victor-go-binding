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

func InsertVectorHandler(w http.ResponseWriter, r *http.Request) {
	var req models.InsertVectorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Error decoding request:", err)
		http.Error(w, "Invalid insert vector request payload", http.StatusBadRequest)
		return
	}

	urlParams := mux.Vars(r)
	indexNameParam := urlParams["indexName"]

	indexResource, exists := store.GetIndex(indexNameParam)
	if !exists {
		http.Error(w, "Index not found", http.StatusNotFound)
		return
	}
	vIndex := indexResource.VIndex

	dims, dimsExists := store.GetIndexDims(indexNameParam)
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

	response := models.InsertVectorResponse{
		Status:  "success",
		Message: "Vector inserted successfully",
		Results: models.InsertVectorResult{
			ID:     req.ID,
			Vector: req.Vector,
		},
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func DeleteVectorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--DELETE HANDLER--")
	vars := mux.Vars(r)
	indexName := vars["indexName"]
	vectorID := vars["vectorID"]

	indexResource, exists := store.GetIndex(indexName)
	if !exists {
		http.Error(w, "Index not found", http.StatusNotFound)
		return
	}

	vIndex := indexResource.VIndex

	id, err := strconv.Atoi(vectorID)
	if err != nil {
		http.Error(w, "Invalid vector ID", http.StatusBadRequest)
		return
	}

	err = vIndex.Delete(uint64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.DeleteVectorResponse{
		Status:  "200",
		Message: "Vector deleted successfully",
		Results: models.DeleteVectorResult{
			ID: uint64(id),
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SearchVectorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--SEARCH HANDLER--")
	vars := mux.Vars(r)
	indexName := vars["indexName"]

	indexResource, exists := store.GetIndex(indexName)
	if !exists {
		http.Error(w, "Index not found", http.StatusNotFound)
		return
	}

	fmt.Println("Index name:", indexName)

	vIndex := indexResource.VIndex

	vectorParam := r.URL.Query().Get("vector")
	if vectorParam == "" {
		http.Error(w, "Missing vector parameter", http.StatusBadRequest)
		return
	}

	kParam := r.URL.Query().Get("top_k")
	var k int
	var err error
	if kParam != "" {
		k, err = strconv.Atoi(kParam)
		if err != nil {
			http.Error(w, "Invalid k parameter", http.StatusBadRequest)
			return
		}
		if k <= 0 {
			http.Error(w, "k must be greater than 0", http.StatusBadRequest)
			return
		}
	} else {
		k = 5
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

	fmt.Println("Vector to search:", vector)
	fmt.Println("topK:", k)

	result, err := vIndex.Search(vector, k)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.SearchVectorResponse{
		Status:  "success",
		Message: "Search completed successfully",
		Results: models.SearchVectorResult{
			ID:       uint64(result.ID),
			Distance: result.Distance,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
