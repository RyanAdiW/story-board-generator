package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"story-board-generator/internal/adapters/postgres"
	queueadapter "story-board-generator/internal/adapters/redis"
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

	generationService := app.NewGenerationService(repo)
	processor := worker.NewProcessor(generationService)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("worker consuming queue=%s", cfg.RabbitMQQueue)
	if err := consumer.ConsumeStoryboardGenerate(ctx, processor.ProcessStoryboardGenerate); err != nil {
		log.Fatalf("consume messages: %v", err)
	}
}
