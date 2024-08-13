package config

type Config struct {
	RunAddr string
	BaseURL string
}

func Load(runAddr string, baseURL string) *Config {
	config := &Config{
		BaseURL: runAddr,
		RunAddr: baseURL,
	}
	return config
}
