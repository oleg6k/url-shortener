package app

import "math/rand"

type Service struct {
	storage map[string]string
}

func NewService() *Service {
	return &Service{storage: make(map[string]string)}
}

func (service *Service) getHashByUrl(url string) string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	hash := string(b)

	service.storage[hash] = url
	if service.storage[url] != "" {
		delete(service.storage, url)
	}
	service.storage[url] = hash

	return hash
}
func (service *Service) getUrlByHash(hash string) string {
	return service.storage[hash]
}
