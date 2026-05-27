package ports

import (
	"context"

	"story-board-generator/internal/domain"
)

type ImageClient interface {
	GenerateSceneImage(ctx context.Context, input SceneImageInput) (SceneImageOutput, error)
}

type SceneImageInput struct {
	Scene           domain.Scene
	Project         domain.Project
	ReferenceAssets []domain.Asset
}

type SceneImageOutput struct {
	Bytes    []byte
	MimeType string
}
