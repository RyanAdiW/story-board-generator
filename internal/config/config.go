package config

import "os"

type Config struct {
	Port               string
	UploadDir          string
	DataDir            string
	RabbitMQURL        string
	RabbitMQQueue      string
	OpenAIAPIKey       string
	OpenAITextModel    string
	OpenAIImageModel   string
	OpenAIImageSize    string
	OpenAIImageQuality string
}

func FromEnv() Config {
	return Config{
		Port:               getenv("APP_PORT", "8080"),
		UploadDir:          getenv("UPLOAD_DIR", "./uploads"),
		DataDir:            getenv("DATA_DIR", "./data"),
		RabbitMQURL:        getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitMQQueue:      getenv("RABBITMQ_QUEUE", "storyboard.generate"),
		OpenAIAPIKey:       getenv("OPENAI_API_KEY", ""),
		OpenAITextModel:    getenv("OPENAI_TEXT_MODEL", "gpt-4.1-mini"),
		OpenAIImageModel:   getenv("OPENAI_IMAGE_MODEL", "gpt-image-1"),
		OpenAIImageSize:    getenv("OPENAI_IMAGE_SIZE", "1024x1536"),
		OpenAIImageQuality: getenv("OPENAI_IMAGE_QUALITY", "medium"),
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
