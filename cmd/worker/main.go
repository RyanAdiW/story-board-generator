package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"story-board-generator/internal/config"
	"story-board-generator/internal/queue"
	"story-board-generator/internal/store"
	"story-board-generator/internal/worker"
)

func main() {
	cfg := config.FromEnv()

	repo, err := store.New(cfg.DataDir)
	if err != nil {
		log.Fatalf("init store: %v", err)
	}

	consumer, err := queue.NewConsumer(cfg.RabbitMQURL, cfg.RabbitMQQueue)
	if err != nil {
		log.Fatalf("init rabbitmq consumer: %v", err)
	}
	defer consumer.Close()

	processor := worker.NewProcessor(repo)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Printf("worker consuming queue=%s", cfg.RabbitMQQueue)
	if err := consumer.ConsumeStoryboardGenerate(ctx, processor.ProcessStoryboardGenerate); err != nil {
		log.Fatalf("consume messages: %v", err)
	}
}
