package ports

import (
	"context"

	"story-board-generator/internal/domain"
)

type Renderer interface {
	RenderStoryboard(ctx context.Context, input RenderStoryboardInput) (RenderStoryboardOutput, error)
}

type RenderStoryboardInput struct {
	Project domain.Project
	Scenes  []domain.Scene
}

type RenderStoryboardOutput struct {
	Bytes    []byte
	MimeType string
}
