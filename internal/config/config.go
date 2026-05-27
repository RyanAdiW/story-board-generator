package config

import "os"

type Config struct {
	Port      string
	UploadDir string
	DataDir   string
}

func FromEnv() Config {
	return Config{
		Port:      getenv("APP_PORT", "8080"),
		UploadDir: getenv("UPLOAD_DIR", "./uploads"),
		DataDir:   getenv("DATA_DIR", "./data"),
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
