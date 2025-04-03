package store

import (
	"sync"
	"victorgo/daemon/pkg/models"
)

var (
	indexStore = make(map[string]*models.IndexResource)
	indexMutex sync.RWMutex
)

func GetIndex(indexName string) (*models.IndexResource, bool) {
	indexMutex.RLock()
	defer indexMutex.RUnlock()

	index, exists := indexStore[indexName]
	return index, exists
}

func GetIndexDims(indexName string) (uint16, bool) {
	indexMutex.RLock()
	defer indexMutex.RUnlock()

	index, exists := indexStore[indexName]
	if !exists {
		return 0, false
	}
	return index.Dims, true
}

func StoreIndex(indexResource *models.IndexResource) {
	indexMutex.Lock()
	defer indexMutex.Unlock()

	indexStore[indexResource.IndexName] = indexResource
}
