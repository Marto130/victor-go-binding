package victorSDK

import (
	"net/http"
)

type Client struct {
	HttpClient *http.Client
	BaseURL    string
	IsLocal    bool
}

// VB = Victor binding
type VBMatchResult struct {
	ID       int     `json:"id"`
	Distance float32 `json:"distance"`
}

type ClientOptions struct {
	Host            string
	Port            string
	AutoStartDaemon bool
}

type CommandOutput struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Results interface{} `json:"results"`
}

type CreateIndexCommandInput struct {
	IndexType int    `json:"index_type"`
	Method    int    `json:"method"`
	Dims      uint16 `json:"dims"`
	IndexName string `json:"index_name"`
}

type CreateIndexCommandOutput struct {
	CommandOutput
	Results struct {
		IndexName string `json:"index_name"`
		ID        string `json:"id"`
		Dims      uint16 `json:"dims"`
		IndexType int    `json:"index_type"`
		Method    int    `json:"method"`
	} `json:"results"`
}

type InsertVectorCommandInput struct {
	IndexName string    `json:"index_name"`
	ID        uint64    `json:"id"`
	Vector    []float32 `json:"vector"`
}

type InsertVectorCommandOutput struct {
	CommandOutput
	Results struct {
		ID     uint64    `json:"id"`
		Vector []float32 `json:"vector"`
	} `json:"results"`
}

type DeleteVectorCommandInput struct {
	IndexName string `json:"index_name"`
	VectorID  uint64 `json:"vector_id"`
}

type DeleteVectorCommandOutput struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Results struct {
		ID uint64 `json:"id"`
	} `json:"results"`
}

type SearchVectorCommandInput struct {
	IndexName string    `json:"index_name"`
	TopK      int       `json:"top_k"`
	Vector    []float32 `json:"vector"`
}

type SearchCommandOutput struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Results VBMatchResult `json:"results"`
}
