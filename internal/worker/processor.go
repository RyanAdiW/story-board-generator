package worker

import (
	"context"
	"fmt"
	"time"

	"story-board-generator/internal/queue"
	"story-board-generator/internal/store"
)

type Processor struct {
	repo *store.Store
}

func NewProcessor(repo *store.Store) *Processor {
	return &Processor{repo: repo}
}

func (p *Processor) ProcessStoryboardGenerate(_ context.Context, payload queue.StoryboardGeneratePayload) error {
	if err := p.repo.UpdateJobProcessing(payload.ProjectID, payload.JobID, "analyzing_product"); err != nil {
		return fmt.Errorf("set processing state: %w", err)
	}
	time.Sleep(1 * time.Second)

	if err := p.repo.UpdateJobProcessing(payload.ProjectID, payload.JobID, "generating_scenes"); err != nil {
		_ = p.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scenes", err.Error())
		return fmt.Errorf("update scene generation state: %w", err)
	}
	time.Sleep(1 * time.Second)

	if err := p.repo.UpdateJobCompleted(payload.ProjectID, payload.JobID); err != nil {
		_ = p.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "finalizing", err.Error())
		return fmt.Errorf("set completed state: %w", err)
	}

	return nil
}
