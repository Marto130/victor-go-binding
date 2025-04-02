package db

import (
	"fmt"
	victor "victorgo/binding"
)

// VectorDB represents the vector database interface
type VectorDB struct {
	index *victor.Index
}

// NewVectorDB initializes a new VectorDB instance
func NewVectorDB(indexType, method int, dims uint16) (*VectorDB, error) {
	index, err := victor.AllocIndex(indexType, method, dims)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate index: %w", err)
	}
	return &VectorDB{index: index}, nil
}

// Insert inserts a vector into the database
func (db *VectorDB) Insert(id uint64, vector []float32) error {
	return db.index.Insert(id, vector)
}

// Search searches for the closest match for a given vector
func (db *VectorDB) Search(vector []float32, dims int) (*victor.MatchResult, error) {
	return db.index.Search(vector, dims)
}

// Delete removes a vector from the database by its ID
func (db *VectorDB) Delete(id uint64) error {
	return db.index.Delete(id)
}

// Close releases the resources associated with the VectorDB
func (db *VectorDB) Close() {
	db.index.DestroyIndex()
}