package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	aiadapter "story-board-generator/internal/adapters/ai"
	"story-board-generator/internal/adapters/postgres"
	queueadapter "story-board-generator/internal/adapters/rabbitmq"
	"story-board-generator/internal/adapters/renderer"
	"story-board-generator/internal/adapters/storage"
	"story-board-generator/internal/app"
	"story-board-generator/internal/config"
	"story-board-generator/internal/worker"
)

func main() {
	cfg := config.FromEnv()

	repo, err := postgres.NewRepository(cfg.DataDir)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	consumer, err := queueadapter.NewQueue(cfg.RabbitMQURL, cfg.RabbitMQQueue)
	if err != nil {
		log.Fatalf("init rabbitmq consumer: %v", err)
	}
	defer consumer.Close()

	sceneAIClient := aiadapter.NewOpenAIClient(cfg.OpenAIAPIKey, cfg.OpenAITextModel)
	imageAIClient := aiadapter.NewImageClient(cfg.OpenAIAPIKey, cfg.OpenAIImageModel, cfg.OpenAIImageSize, cfg.OpenAIImageQuality)
	generationService := app.NewGenerationService(
		repo,
		sceneAIClient,
		imageAIClient,
		renderer.NewStoryboardRenderer(),
		storage.NewS3Storage(cfg.UploadDir),
	)
	processor := worker.NewProcessor(generationService)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("worker consuming queue=%s", cfg.RabbitMQQueue)
	if err := consumer.ConsumeStoryboardGenerate(ctx, processor.ProcessStoryboardGenerate); err != nil {
		log.Fatalf("consume messages: %v", err)
	}
}
