package app

import (
	"fmt"
	"github.com/oleg6k/url-shortener/internal/app/config"
	"github.com/oleg6k/url-shortener/internal/app/repositories"
	"github.com/oleg6k/url-shortener/internal/app/types"
)

type Storage struct {
	repository types.RepositoryInterface
}

func NewStorage(storageCfg config.StorageConfig) (*Storage, error) {
	var repo types.RepositoryInterface
	var err error

	if storageCfg.Database.URL != "" {
		repo, err = repositories.NewDatabaseRepository(storageCfg.Database)
		fmt.Println("Database repository...")
	} else if storageCfg.Disk.Path != "" {
		repo, err = repositories.NewDiskRepository(storageCfg.Disk)
		fmt.Println("Disk repository...")
	} else {
		repo, err = repositories.NewInMemoryRepository()
		fmt.Println("Memory repository...")
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
