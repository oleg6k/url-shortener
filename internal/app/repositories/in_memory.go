package repositories

import (
	"github.com/oleg6k/url-shortener/internal/app/types"
	"sync"
)

type InMemoryRepository struct {
	sync.RWMutex
	storage map[string]types.URLRecord
}

func NewInMemoryRepository() (*InMemoryRepository, error) {
	return &InMemoryRepository{
		storage: make(map[string]types.URLRecord),
	}, nil
}

func (r *InMemoryRepository) Add(record types.URLRecord) error {
	r.Lock()
	defer r.Unlock()

	r.storage[record.ShortURL] = record
	r.storage[record.OriginalURL] = record

	return nil
}

func (r *InMemoryRepository) Get(key string) (types.URLRecord, bool) {
	r.RLock()
	defer r.RUnlock()
	value, exists := r.storage[key]
	return value, exists
}

func (r *InMemoryRepository) Delete(key string) error {
	r.Lock()
	defer r.Unlock()
	delete(r.storage, key)
	return nil
}

func (r *InMemoryRepository) Health() error {
	return nil
}

func (r *InMemoryRepository) Close() error {
	return nil
}
