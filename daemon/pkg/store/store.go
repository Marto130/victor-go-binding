package store

import (
	"sync"
	"victorgo/daemon/pkg/models"
)

var (
	indexStore = make(map[string]*models.IndexResource)
	indexMutex sync.RWMutex
)

func GetIndex(indexID string) (*models.IndexResource, bool) {
	indexMutex.RLock()
	defer indexMutex.RUnlock()

	index, exists := indexStore[indexID]
	return index, exists
}

func GetIndexDims(indexID string) (uint16, bool) {
	indexMutex.RLock()
	defer indexMutex.RUnlock()

	index, exists := indexStore[indexID]
	if !exists {
		return 0, false
	}
	return index.Dims, true
}

func StoreIndex(indexResource *models.IndexResource) {
	indexMutex.Lock()
	defer indexMutex.Unlock()

	indexStore[indexResource.IndexID] = indexResource
}
