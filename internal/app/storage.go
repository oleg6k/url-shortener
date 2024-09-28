package app

import (
	"github.com/oleg6k/url-shortener/internal/app/repositories"
	"github.com/oleg6k/url-shortener/internal/app/types"
)

type Storage struct {
	repository types.RepositoryInterface
}

func NewStorage(databaseURL, diskPath string) (*Storage, error) {
	var repo types.RepositoryInterface
	var err error

	if databaseURL != "" {
		repo, err = repositories.NewDatabaseRepository(databaseURL)
	} else if diskPath != "" {
		repo, err = repositories.NewDiskRepository(diskPath)
	} else {
		repo, err = repositories.NewInMemoryRepository()
	}

	if err != nil {
		return nil, err
	}

	return &Storage{repository: repo}, nil
}

func (p *Storage) Add(record types.URLRecord) error {
	return p.repository.Add(record)
}

func (p *Storage) Get(key string) (types.URLRecord, bool) {
	return p.repository.Get(key)
}

func (p *Storage) Delete(key string) error {
	return p.repository.Delete(key)
}

func (p *Storage) Health() error {
	return p.repository.Health()
}

func (p *Storage) Close() error {
	return p.repository.Close()
}
