package config

import (
	"flag"
	"os"
)

type AppConfig struct {
	RunAddr string
	BaseURL string
}

type DatabaseConfig struct {
	URL           string
	MigrationsDir string
}

type DiskConfig struct {
	Path string
}

type StorageConfig struct {
	Database DatabaseConfig
	Disk     DiskConfig
}

type Config struct {
	App     AppConfig
	Storage StorageConfig
}

func Load() *Config {

	var runAddr string
	var baseURL string
	var fileStoragePath string
	var databaseURL string

	flag.StringVar(&runAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base address and port to shortener results")
	flag.StringVar(&fileStoragePath, "f", "", "Storage Path")
	flag.StringVar(&databaseURL, "d", "", "Database url")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		runAddr = envRunAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = envBaseURL
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}

	if envDatabaseURL := os.Getenv("DATABASE_DSN"); envDatabaseURL != "" {
		databaseURL = envDatabaseURL
	}

	return &Config{
		App: AppConfig{
			RunAddr: runAddr,
			BaseURL: baseURL,
		},
		Storage: StorageConfig{
			Database: DatabaseConfig{
				URL:           databaseURL,
				MigrationsDir: "internal/app/migrations",
			},
			Disk: DiskConfig{
				Path: fileStoragePath,
			},
		},
	}
}
