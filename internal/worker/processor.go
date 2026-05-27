package worker

import (
	"context"

	"story-board-generator/internal/ports"
)

type StoryboardGenerator interface {
	ProcessStoryboardGenerate(ctx context.Context, payload ports.StoryboardGeneratePayload) error
}

type Processor struct {
	generator StoryboardGenerator
}

func NewProcessor(generator StoryboardGenerator) *Processor {
	return &Processor{generator: generator}
}

func (p *Processor) ProcessStoryboardGenerate(ctx context.Context, payload ports.StoryboardGeneratePayload) error {
	return p.generator.ProcessStoryboardGenerate(ctx, payload)
}
