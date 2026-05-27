package config

import "os"

type Config struct {
	Port            string
	UploadDir       string
	DataDir         string
	RabbitMQURL     string
	RabbitMQQueue   string
	OpenAIAPIKey    string
	OpenAITextModel string
}

func FromEnv() Config {
	return Config{
		Port:            getenv("APP_PORT", "8080"),
		UploadDir:       getenv("UPLOAD_DIR", "./uploads"),
		DataDir:         getenv("DATA_DIR", "./data"),
		RabbitMQURL:     getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitMQQueue:   getenv("RABBITMQ_QUEUE", "storyboard.generate"),
		OpenAIAPIKey:    getenv("OPENAI_API_KEY", ""),
		OpenAITextModel: getenv("OPENAI_TEXT_MODEL", "gpt-4.1-mini"),
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
