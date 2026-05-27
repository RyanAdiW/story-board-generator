package main

import (
	"log"

	"story-board-generator/internal/config"
	httpapi "story-board-generator/internal/http"
	"story-board-generator/internal/queue"
	"story-board-generator/internal/store"
)

func main() {
	cfg := config.FromEnv()

	repo, err := store.New(cfg.DataDir)
	if err != nil {
		log.Fatalf("init store: %v", err)
	}

	queueClient, err := queue.NewClient(cfg.RabbitMQURL, cfg.RabbitMQQueue)
	if err != nil {
		log.Fatalf("init queue client: %v", err)
	}
	defer queueClient.Close()

	handler := httpapi.NewHandler(repo, cfg.UploadDir, queueClient)
	e := httpapi.NewRouter(handler)

	log.Printf("api listening on :%s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("start api: %v", err)
	}
}
