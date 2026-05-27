package app

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"story-board-generator/internal/domain"
	"story-board-generator/internal/ports"
)

type GenerationService struct {
	repo    ports.Repository
	sceneAI ports.AIClient
	imageAI ports.ImageClient
	storage ports.Storage
}

func NewGenerationService(repo ports.Repository, sceneAI ports.AIClient, imageAI ports.ImageClient, storage ports.Storage) *GenerationService {
	return &GenerationService{
		repo:    repo,
		sceneAI: sceneAI,
		imageAI: imageAI,
		storage: storage,
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

	scenes, err := s.sceneAI.GenerateScenes(ctx, ports.SceneGenerationInput{
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

	if err := s.repo.UpdateJobProcessing(payload.ProjectID, payload.JobID, "generating_scene_images"); err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scene_images", err.Error())
		return fmt.Errorf("update scene image generation state: %w", err)
	}

	updatedScenes, generatedAssets, err := s.generateSceneAssets(ctx, bundle.Project, bundle.Assets, validScenes)
	if err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scene_images", err.Error())
		return fmt.Errorf("generate scene assets: %w", err)
	}

	if err := s.repo.AddAssets(payload.ProjectID, generatedAssets); err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scene_images", err.Error())
		return fmt.Errorf("save generated assets: %w", err)
	}

	if err := s.repo.SaveScenes(payload.ProjectID, updatedScenes); err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "generating_scene_images", err.Error())
		return fmt.Errorf("save updated scenes: %w", err)
	}

	if err := s.repo.UpdateJobCompleted(payload.ProjectID, payload.JobID); err != nil {
		_ = s.repo.UpdateJobFailed(payload.ProjectID, payload.JobID, "finalizing", err.Error())
		return fmt.Errorf("set completed state: %w", err)
	}

	return nil
}

func (s *GenerationService) generateSceneAssets(ctx context.Context, project domain.Project, referenceAssets []domain.Asset, scenes []domain.Scene) ([]domain.Scene, []domain.Asset, error) {
	now := time.Now().UTC()
	updatedScenes := make([]domain.Scene, 0, len(scenes))
	generatedAssets := make([]domain.Asset, 0, len(scenes))

	for _, scene := range scenes {
		imageResult, err := s.imageAI.GenerateSceneImage(ctx, ports.SceneImageInput{
			Scene:           scene,
			Project:         project,
			ReferenceAssets: referenceAssets,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("scene %d image generation: %w", scene.SceneNumber, err)
		}
		if len(imageResult.Bytes) == 0 {
			return nil, nil, fmt.Errorf("scene %d returned empty image", scene.SceneNumber)
		}

		assetID, err := newID()
		if err != nil {
			return nil, nil, fmt.Errorf("generate asset id: %w", err)
		}
		ext := extensionFromMime(imageResult.MimeType)
		objectPath := filepath.ToSlash(filepath.Join(project.ID, "scenes", assetID+ext))

		fileURL, err := s.storage.Upload(ctx, objectPath, imageResult.MimeType, bytes.NewReader(imageResult.Bytes))
		if err != nil {
			return nil, nil, fmt.Errorf("upload scene %d image: %w", scene.SceneNumber, err)
		}

		asset := domain.Asset{
			ID:        assetID,
			ProjectID: project.ID,
			AssetType: "generated_scene_image",
			FileURL:   fileURL,
			MimeType:  imageResult.MimeType,
			CreatedAt: now,
		}
		scene.ImageAssetID = assetID
		scene.ImageURL = fileURL
		updatedScenes = append(updatedScenes, scene)
		generatedAssets = append(generatedAssets, asset)
	}

	return updatedScenes, generatedAssets, nil
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

func extensionFromMime(mimeType string) string {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	default:
		return ".png"
	}
}
