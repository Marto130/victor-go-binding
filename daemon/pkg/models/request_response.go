package models

type CreateIndexResult struct {
	IndexName string `json:"index_name"`
	ID        string `json:"id"`
	Dims      uint16 `json:"dims"`
	IndexType int    `json:"index_type"`
	Method    int    `json:"method"`
}

type InsertVectorResult struct {
	ID     uint64    `json:"id"`
	Vector []float32 `json:"vector"`
}

type DeleteVectorResult struct {
	ID uint64 `json:"id"`
}

type SearchVectorResult struct {
	ID       uint64  `json:"id"`
	Distance float32 `json:"distance"`
}

type CreateIndexRequest struct {
	IndexType int    `json:"index_type"`
	Method    int    `json:"method"`
	Dims      uint16 `json:"dims"`
}

type CreateIndexResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message,omitempty"`
	Results CreateIndexResult `json:"results"`
}

type InsertVectorRequest struct {
	ID     uint64    `json:"id"`
	Vector []float32 `json:"vector"`
}

type InsertVectorResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message,omitempty"`
	Results InsertVectorResult `json:"results"`
}

type DeleteVectorResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message,omitempty"`
	Results DeleteVectorResult `json:"results"`
}

type SearchVectorRequest struct {
	TopK   int       `json:"top_k"`
	Vector []float32 `json:"vector"`
}
type SearchVectorResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message,omitempty"`
	Results SearchVectorResult `json:"results"`
}
