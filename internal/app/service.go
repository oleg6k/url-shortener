package app

import (
	"github.com/oleg6k/url-shortener/internal/app/types"
	"math/rand"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}

func (service *Service) getHashByURL(url string) (string, error) {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	hash := string(b)

	err := service.storage.Add(types.URLRecord{
		ShortURL:    hash,
		OriginalURL: url,
	})

	if err != nil {
		return "", err
	}

	return hash, nil
}

func (service *Service) getURLByHash(hash string) string {
	url, _ := service.storage.Get(hash)
	return url.OriginalURL
}
