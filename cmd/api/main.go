package main

import (
	"log"

	"story-board-generator/internal/config"
	httpapi "story-board-generator/internal/http"
	"story-board-generator/internal/store"
)

func main() {
	cfg := config.FromEnv()

	repo, err := store.New(cfg.DataDir)
	if err != nil {
		log.Fatalf("init store: %v", err)
	}

	handler := httpapi.NewHandler(repo, cfg.UploadDir)
	e := httpapi.NewRouter(handler)

	log.Printf("api listening on :%s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("start api: %v", err)
	}
}
