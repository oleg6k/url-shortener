package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr         string
	BaseURL         string
	FileStoragePath string
	DatabaseURL     string
}

func Load() *Config {

	var runAddr string
	var baseURL string
	var fileStoragePath string
	var databaseURL string

	flag.StringVar(&runAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base address and port to shortener results")
	flag.StringVar(&fileStoragePath, "f", "./tmp/storage.json", "storage path")
	flag.StringVar(&databaseURL, "d", "postgresql://postgres:postgres@localhost:5432/postgres", "database url")
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
		RunAddr:         runAddr,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
		DatabaseURL:     databaseURL,
	}
}
