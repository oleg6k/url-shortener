package repositories

import (
	"bufio"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/oleg6k/url-shortener/internal/app/config"
	"github.com/oleg6k/url-shortener/internal/app/types"
	"os"
)

type DiskRepository struct {
	filePath           string
	inMemoryRepository *InMemoryRepository
}

func NewDiskRepository(cfg config.DiskConfig) (*DiskRepository, error) {
	inMemoryRepository, err := NewInMemoryRepository()
	if err != nil {
		return nil, err
	}
	repository := &DiskRepository{filePath: cfg.Path,
		inMemoryRepository: inMemoryRepository}

	err = repository.Load()
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *DiskRepository) Add(record types.URLRecord) error {

	if _, exist := r.inMemoryRepository.Get(record.OriginalURL); exist {
		return nil
	}

	fullRecord := types.DiskURLRecord{
		UUID:      uuid.New().String(),
		URLRecord: record,
	}

	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(fullRecord)
	if err != nil {
		return err
	}

	if _, err = file.Write(append(data, '\n')); err != nil {
		return err
	}

	return r.inMemoryRepository.Add(record)
}

func (r *DiskRepository) Get(key string) (types.URLRecord, bool) {
	return r.inMemoryRepository.Get(key)
}

func (r *DiskRepository) Delete(key string) error {
	return r.inMemoryRepository.Delete(key)
}

func (r *DiskRepository) Load() error {

	file, err := os.OpenFile(r.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record types.URLRecord
		if err = json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return err
		}
		if err = r.inMemoryRepository.Add(record); err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (r *DiskRepository) Health() error {
	return nil
}

func (r *DiskRepository) Close() error {
	return nil
}
