package app

import (
	"bufio"
	"encoding/json"
	"github.com/google/uuid"
	"os"
	"sync"
)

type URLRecord struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Repository struct {
	sync.RWMutex
	storage  map[string]string
	filePath string
}

func NewRepository(filePath string) *Repository {
	return &Repository{filePath: filePath, storage: make(map[string]string)}
}

func (r *Repository) Load() error {
	r.Lock()
	defer r.Unlock()

	file, err := os.OpenFile(r.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record URLRecord
		if err = json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return err
		}
		r.storage[record.ShortURL] = record.OriginalURL
		r.storage[record.OriginalURL] = record.ShortURL
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) AddURLRecord(record URLRecord) error {

	if r.storage[record.OriginalURL] != "" {
		return nil
	}

	r.Lock()
	defer r.Unlock()

	recordUUID := uuid.New().String()
	fullRecord := struct {
		UUID string `json:"uuid"`
		URLRecord
	}{
		UUID:      recordUUID,
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

	r.storage[record.OriginalURL] = record.ShortURL
	r.storage[record.ShortURL] = record.OriginalURL
	return nil
}
