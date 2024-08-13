package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr string
	BaseURL string
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func Load() *Config {
	runAddr := getEnv("RUN_ADDRESS", ":8080")
	baseURL := getEnv("BASE_URL", "http://localhost:8080")

	flag.StringVar(&runAddr, "a", runAddr, "address and port to run server")
	flag.StringVar(&baseURL, "b", baseURL, "base address and port to shortener results")
	flag.Parse()

	return &Config{
		RunAddr: runAddr,
		BaseURL: baseURL,
	}
}
