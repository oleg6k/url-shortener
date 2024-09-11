package app

import "math/rand"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (service *Service) getHashByURL(url string) (string, error) {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	hash := string(b)

	err := service.repository.AddURLRecord(URLRecord{
		ShortURL:    hash,
		OriginalURL: url,
	})

	if err != nil {
		return "", err
	}

	return service.repository.storage[url], nil
}

func (service *Service) getURLByHash(hash string) string {
	return service.repository.storage[hash]
}
