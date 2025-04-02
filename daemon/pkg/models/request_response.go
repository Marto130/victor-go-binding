package models

import (
	binding "victorgo/binding"
)

// Request represents the structure of a request to the API.
type Request struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

// Response represents the structure of a response from the API.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// IndexRequest represents a request to create an index
type CreateIndexRequest struct {
	IndexType int    `json:"index_type"`
	Method    int    `json:"method"`
	Dims      uint16 `json:"dims"`
}

// InsertVectorRequest represents a request to insert a vector into the database
type InsertVectorRequest struct {
	ID     uint64    `json:"id"`
	Vector []float32 `json:"vector"`
}

// SearchVectorRequest represents a request to search for a vector in the database
// type SearchVectorRequest struct {
// 	Vector  []float32 `json:"vector"`
// 	IndexID string    `json:"index_id"`
// }

type IndexResource struct {
	CreateIndexRequest
	IndexID string `json:"index_id"`
	VIndex  *binding.Index
}
