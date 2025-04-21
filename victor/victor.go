package victor

/*
#cgo CFLAGS: -I/Users/juan.irigoin/workspace/victor/victorgo/include
#cgo LDFLAGS: -L/Users/juan.irigoin/workspace/victor/victorgo/lib -lvictor
#include <victor/victor.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ErrorCode maps C error codes to Go
type ErrorCode int

const (
	SUCCESS ErrorCode = iota
	INVALID_INIT
	INVALID_INDEX
	INVALID_VECTOR
	INVALID_RESULT
	INVALID_DIMENSIONS
	INVALID_ARGUMENT
	INVALID_ID
	INVALID_REF
	DUPLICATED_ENTRY
	NOT_FOUND_ID
	INDEX_EMPTY
	THREAD_ERROR
	SYSTEM_ERROR
	NOT_IMPLEMENTED
)

// errorMessages maps error codes to human-readable messages
var errorMessages = map[ErrorCode]string{
	SUCCESS:            "Success",
	INVALID_INIT:       "Invalid initialization",
	INVALID_INDEX:      "Invalid index",
	INVALID_VECTOR:     "Invalid vector",
	INVALID_RESULT:     "Invalid result",
	INVALID_DIMENSIONS: "Invalid dimensions",
	INVALID_ARGUMENT:   "Invalid argument",
	INVALID_ID:         "Invalid ID",
	INVALID_REF:        "Invalid reference",
	DUPLICATED_ENTRY:   "Duplicated entry",
	NOT_FOUND_ID:       "ID not found",
	INDEX_EMPTY:        "Index is empty",
	THREAD_ERROR:       "Thread error",
	SYSTEM_ERROR:       "System error",
	NOT_IMPLEMENTED:    "Not implemented",
}

// toError converts a C error code to a Go error
func toError(code C.int) error {
	if code == C.int(SUCCESS) {
		return nil
	}
	if msg, exists := errorMessages[ErrorCode(code)]; exists {
		return fmt.Errorf(msg)
	}
	return fmt.Errorf("Unknown error code: %d", code)
}

// MatchResult represents a search result in Go
type MatchResult struct {
	ID       int     `json:"id"`
	Distance float32 `json:"distance"`
}

// Index represents an index structure in Go
type Index struct {
	ptr *C.Index
}

// AllocIndex creates a new index
func AllocIndex(indexType, method int, dims uint16) (*Index, error) {
	idx := C.alloc_index(C.int(indexType), C.int(method), C.uint16_t(dims), nil)
	if idx == nil {
		return nil, fmt.Errorf("Failed to allocate index")
	}
	return &Index{ptr: idx}, nil
}

// Insert adds a vector to the index with a given ID
func (idx *Index) Insert(id uint64, vector []float32) error {
	if idx.ptr == nil {
		return fmt.Errorf("Index not initialized")
	}
	if len(vector) == 0 {
		return fmt.Errorf("Empty vector")
	}

	cVector := (*C.float)(unsafe.Pointer(&vector[0]))
	return toError(C.insert(idx.ptr, C.uint64_t(id), cVector, C.uint16_t(len(vector))))
}

// Search finds the closest match for a given vector
func (idx *Index) Search(vector []float32, dims int) (*MatchResult, error) {
	if idx.ptr == nil {
		return nil, fmt.Errorf("Index not initialized")
	}

	var cResult C.MatchResult
	cVector := (*C.float)(unsafe.Pointer(&vector[0]))
	err := C.search(idx.ptr, cVector, C.uint16_t(dims), &cResult)
	if e := toError(err); e != nil {
		return nil, e
	}

	return &MatchResult{
		ID:       int(cResult.id),
		Distance: float32(cResult.distance),
	}, nil
}

// SearchN finds the n closest matches for a given vector
func (idx *Index) SearchN(vector []float32, n int) ([]MatchResult, error) {
	if idx.ptr == nil {
		return nil, fmt.Errorf("Index not initialized")
	}

	if len(vector) == 0 {
		return nil, fmt.Errorf("Empty vector")
	}

	// Allocate memory for results
	cResults := make([]C.MatchResult, n)
	cVector := (*C.float)(unsafe.Pointer(&vector[0]))

	// Call the C function
	err := C.search_n(idx.ptr, cVector, C.uint16_t(len(vector)),
		(*C.MatchResult)(unsafe.Pointer(&cResults[0])), C.int(n))

	if e := toError(err); e != nil {
		return nil, e
	}

	// Convert C results to Go results
	results := make([]MatchResult, n)
	for i := 0; i < n; i++ {
		results[i] = MatchResult{
			ID:       int(cResults[i].id),
			Distance: float32(cResults[i].distance),
		}
	}

	return results, nil
}

// Delete removes a vector from the index by its ID
func (idx *Index) Delete(id uint64) error {
	if idx.ptr == nil {
		return fmt.Errorf("Index not initialized")
	}
	return toError(C.delete(idx.ptr, C.uint64_t(id)))
}

// DestroyIndex releases index memory
func (idx *Index) DestroyIndex() {
	if idx.ptr != nil {
		C.destroy_index(&idx.ptr)
		idx.ptr = nil
	}
}

// TimeStat represents timing statistics for an operation
type TimeStat struct {
	Count uint64  `json:"count"` // Number of operations
	Total float64 `json:"total"` // Total time in seconds
	Last  float64 `json:"last"`  // Last operation time
	Min   float64 `json:"min"`   // Minimum operation time
	Max   float64 `json:"max"`   // Maximum operation time
}

// IndexStats represents aggregate statistics for the index
type IndexStats struct {
	Insert  TimeStat `json:"insert"`   // Insert operations timing
	Delete  TimeStat `json:"delete"`   // Delete operations timing
	Dump    TimeStat `json:"dump"`     // Dump to file operation
	Search  TimeStat `json:"search"`   // Single search timing
	SearchN TimeStat `json:"search_n"` // Multi-search timing
}

// GetStats retrieves the internal statistics of the index
func (idx *Index) GetStats() (*IndexStats, error) {
	if idx.ptr == nil {
		return nil, fmt.Errorf("Index not initialized")
	}

	// Create a C structure to hold the stats
	var cStats C.IndexStats

	// Call the C function
	err := C.stats(idx.ptr, &cStats)

	if e := toError(err); e != nil {
		return nil, e
	}

	// Convert C stats to Go stats
	stats := &IndexStats{
		Insert: TimeStat{
			Count: uint64(cStats.insert.count),
			Total: float64(cStats.insert.total),
			Last:  float64(cStats.insert.last),
			Min:   float64(cStats.insert.min),
			Max:   float64(cStats.insert.max),
		},
		Delete: TimeStat{
			Count: uint64(cStats.delete.count),
			Total: float64(cStats.delete.total),
			Last:  float64(cStats.delete.last),
			Min:   float64(cStats.delete.min),
			Max:   float64(cStats.delete.max),
		},
		Dump: TimeStat{
			Count: uint64(cStats.dump.count),
			Total: float64(cStats.dump.total),
			Last:  float64(cStats.dump.last),
			Min:   float64(cStats.dump.min),
			Max:   float64(cStats.dump.max),
		},
		Search: TimeStat{
			Count: uint64(cStats.search.count),
			Total: float64(cStats.search.total),
			Last:  float64(cStats.search.last),
			Min:   float64(cStats.search.min),
			Max:   float64(cStats.search.max),
		},
		SearchN: TimeStat{
			Count: uint64(cStats.search_n.count),
			Total: float64(cStats.search_n.total),
			Last:  float64(cStats.search_n.last),
			Min:   float64(cStats.search_n.min),
			Max:   float64(cStats.search_n.max),
		},
	}

	return stats, nil
}

// Size returns the current number of elements in the index
func (idx *Index) Size() (uint64, error) {
	if idx.ptr == nil {
		return 0, fmt.Errorf("Index not initialized")
	}

	// Create a variable to hold the size
	var size C.uint64_t

	// Call the C function
	err := C.size(idx.ptr, &size)

	if e := toError(err); e != nil {
		return 0, e
	}

	// Convert C uint64_t to Go uint64
	return uint64(size), nil
}

// Contains checks whether a given vector ID exists in the index
func (idx *Index) Contains(id uint64) (bool, error) {
	if idx.ptr == nil {
		return false, fmt.Errorf("Index not initialized")
	}

	// Call the C function
	result := C.contains(idx.ptr, C.uint64_t(id))

	// The C function returns 1 if the ID is found, 0 if not
	// Convert this to a Go boolean
	return result == 1, nil
}
