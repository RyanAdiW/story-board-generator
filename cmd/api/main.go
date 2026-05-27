package main

import (
	"log"

	httpadapter "story-board-generator/internal/adapters/http"
	"story-board-generator/internal/adapters/postgres"
	queueadapter "story-board-generator/internal/adapters/redis"
	"story-board-generator/internal/adapters/storage"
	"story-board-generator/internal/app"
	"story-board-generator/internal/config"
)

func main() {
	cfg := config.FromEnv()

	repo, err := postgres.NewRepository(cfg.DataDir)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	queueClient, err := queueadapter.NewQueue(cfg.RabbitMQURL, cfg.RabbitMQQueue)
	if err != nil {
		log.Fatalf("init queue client: %v", err)
	}
	defer queueClient.Close()

	assetService := app.NewAssetService(storage.NewS3Storage(cfg.UploadDir))
	storyboardService := app.NewStoryboardService(repo, queueClient, assetService)

	handler := httpadapter.NewHandler(storyboardService)
	e := httpadapter.NewRouter(handler)

	log.Printf("api listening on :%s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("start api: %v", err)
	}
}
