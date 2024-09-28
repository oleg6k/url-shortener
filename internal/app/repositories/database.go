package repositories

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/oleg6k/url-shortener/internal/app/types"
)

type DatabaseRepository struct {
	DB *sql.DB
}

func NewDatabaseRepository(databaseURL string) (*DatabaseRepository, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		panic(err)
	}
	return &DatabaseRepository{DB: db}, nil
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

func (r *DatabaseRepository) Health() error {
	return r.DB.Ping()
}

func (r *DatabaseRepository) Close() error {
	return r.DB.Close()
}
