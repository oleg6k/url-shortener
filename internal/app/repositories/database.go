package repositories

import (
	"github.com/oleg6k/url-shortener/internal/app/types"
)

type DatabaseRepository struct {
	databaseURL string
}

func NewDatabaseRepository(databaseURL string) (*DatabaseRepository, error) {
	return &DatabaseRepository{databaseURL: databaseURL}, nil
}

func (r *DatabaseRepository) Add(record types.URLRecord) error {
	return nil
}

func (r *DatabaseRepository) Get(key string) (types.URLRecord, bool) {
	return types.URLRecord{}, false
}

func (r *DatabaseRepository) Delete(key string) error {
	return nil
}
