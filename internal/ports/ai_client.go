package ports

import (
	"context"

	"story-board-generator/internal/domain"
)

type AIClient interface {
	GenerateScenes(ctx context.Context, input SceneGenerationInput) ([]domain.Scene, error)
}

type SceneGenerationInput struct {
	Project domain.Project
	Assets  []domain.Asset
}
