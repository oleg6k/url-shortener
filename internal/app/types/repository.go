package types

type RepositoryInterface interface {
	Add(record URLRecord) error
	Get(key string) (URLRecord, bool)
	Delete(key string) error
}

type URLRecord struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type DiskUrlRecord struct {
	UUID string `json:"uuid"`
	URLRecord
}
