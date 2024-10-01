package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/oleg6k/url-shortener/internal/app/config"
	"github.com/oleg6k/url-shortener/internal/app/types"
	"log"
)

type DatabaseRepository struct {
	DB *sql.DB
}

func NewDatabaseRepository(cfg config.DatabaseConfig) (*DatabaseRepository, error) {
	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Connected to the database successfully.")

	if err = runMigrations(db, cfg.MigrationsDir); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DatabaseRepository{DB: db}, nil
}

func runMigrations(db *sql.DB, migrationsDir string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully.")
	return nil
}

func (r *DatabaseRepository) Add(record types.URLRecord) error {
	_, err := r.DB.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", record.ShortURL, record.OriginalURL)
	return err
}

func (r *DatabaseRepository) Get(key string) (types.URLRecord, bool) {
	var record types.URLRecord
	err := r.DB.QueryRow("SELECT short_url, original_url FROM urls WHERE short_url = $1", key).Scan(&record.ShortURL, &record.OriginalURL)
	if err != nil {
		return types.URLRecord{}, false
	}

	return record, true
}

func (r *DatabaseRepository) Delete(key string) error {
	_, err := r.DB.Exec("DELETE FROM urls WHERE short_url=$1", key)
	return err
}

func (r *DatabaseRepository) Health() error {
	return r.DB.Ping()
}

func (r *DatabaseRepository) Close() error {
	return r.DB.Close()
}
