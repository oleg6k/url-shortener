package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr string
	BaseURL string
}

func Load() *Config {
	runAddr := os.Getenv("SERVER_ADDRESS")
	if runAddr == "" {
		flag.StringVar(&runAddr, "a", ":8080", "address and port to run server")
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		flag.StringVar(&baseURL, "b", "http://localhost:8080", "base address and port to shortener results")
	}
	flag.Parse()

	return &Config{
		RunAddr: runAddr,
		BaseURL: baseURL,
	}
}
