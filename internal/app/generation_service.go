package app

import (
	"context"
	"fmt"

	"story-board-generator/internal/domain"
	"story-board-generator/internal/ports"
)

type GenerationService struct {
	repo ports.Repository
	ai   ports.AIClient
}

func NewGenerationService(repo ports.Repository, ai ports.AIClient) *GenerationService {
	return &GenerationService{
		repo: repo,
		ai:   ai,
	}
}

func (s *GenerationService) ProcessStoryboardGenerate(ctx context.Context, payload ports.StoryboardGeneratePayload) error {
	if err := s.repo.UpdateJobProcessing(payload.ProjectID, payload.JobID, "analyzing_product"); err != nil {
		return fmt.Errorf("set processing state: %w", err)
	}

	bundle, err := s.repo.GetProjectBundle(payload.ProjectID)
	if err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "analyzing_product", err.Error())
		return fmt.Errorf("load project bundle: %w", err)
	}

	if err := s.repo.UpdateJobProcessing(payload.ProjectID, payload.JobID, "generating_scenes"); err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scenes", err.Error())
		return fmt.Errorf("update scene generation state: %w", err)
	}

	scenes, err := s.ai.GenerateScenes(ctx, ports.SceneGenerationInput{
		Project: bundle.Project,
		Assets:  bundle.Assets,
	})
	if err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scenes", err.Error())
		return fmt.Errorf("generate scenes: %w", err)
	}

	for i := range scenes {
		scenes[i].ProjectID = payload.ProjectID
		if scenes[i].SceneNumber == 0 {
			scenes[i].SceneNumber = i + 1
		}
	}

	validScenes := normalizeScenes(scenes)
	if len(validScenes) == 0 {
		err = fmt.Errorf("no valid scenes generated")
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scenes", err.Error())
		return err
	}

	if err := s.repo.SaveScenes(payload.ProjectID, validScenes); err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scenes", err.Error())
		return fmt.Errorf("save scenes: %w", err)
	}

	if err := s.repo.UpdateJobCompleted(payload.ProjectID, payload.JobID); err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "finalizing", err.Error())
		return fmt.Errorf("set completed state: %w", err)
	}

	return nil
}

func normalizeScenes(scenes []domain.Scene) []domain.Scene {
	normalized := make([]domain.Scene, 0, len(scenes))
	for _, scene := range scenes {
		if scene.VisualDescription == "" || scene.CameraDirection == "" || scene.MotionDescription == "" || scene.SoundFX == "" || scene.OnScreenText == "" || scene.ImagePrompt == "" {
			continue
		}
		normalized = append(normalized, scene)
	}

	return normalized
}
