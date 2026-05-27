package app

import (
	"context"
	"fmt"
	"time"

	"story-board-generator/internal/ports"
)

type GenerationService struct {
	repo ports.Repository
}

func NewGenerationService(repo ports.Repository) *GenerationService {
	return &GenerationService{repo: repo}
}

func (s *GenerationService) ProcessStoryboardGenerate(_ context.Context, payload ports.StoryboardGeneratePayload) error {
	if err := s.repo.UpdateJobProcessing(payload.ProjectID, payload.JobID, "analyzing_product"); err != nil {
		return fmt.Errorf("set processing state: %w", err)
	}
	time.Sleep(1 * time.Second)

	if err := s.repo.UpdateJobProcessing(payload.ProjectID, payload.JobID, "generating_scenes"); err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scenes", err.Error())
		return fmt.Errorf("update scene generation state: %w", err)
	}
	time.Sleep(1 * time.Second)

	if err := s.repo.UpdateJobCompleted(payload.ProjectID, payload.JobID); err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "finalizing", err.Error())
		return fmt.Errorf("set completed state: %w", err)
	}

	return nil
}
