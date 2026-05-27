package app

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"story-board-generator/internal/domain"
	"story-board-generator/internal/ports"
)

var ErrBadInput = errors.New("bad input")
var ErrNotFound = errors.New("not found")

type StoryboardService struct {
	repo       ports.Repository
	queue      ports.QueuePublisher
	assetStore *AssetService
}

type CreateStoryboardInput struct {
	Title                string
	Style                string
	Platform             string
	Format               string
	TotalDurationSeconds int
	ProductImages        []*multipart.FileHeader
}

type CreateStoryboardOutput struct {
	ProjectID string
	JobID     string
	Status    string
}

func NewStoryboardService(repo ports.Repository, queue ports.QueuePublisher, assetStore *AssetService) *StoryboardService {
	return &StoryboardService{
		repo:       repo,
		queue:      queue,
		assetStore: assetStore,
	}
}

func (s *StoryboardService) CreateStoryboard(ctx context.Context, input CreateStoryboardInput) (CreateStoryboardOutput, error) {
	if strings.TrimSpace(input.Title) == "" ||
		strings.TrimSpace(input.Style) == "" ||
		strings.TrimSpace(input.Platform) == "" ||
		strings.TrimSpace(input.Format) == "" ||
		input.TotalDurationSeconds <= 0 {
		return CreateStoryboardOutput{}, fmt.Errorf("%w: title, style, platform, format, and total_duration_seconds are required", ErrBadInput)
	}

	projectID, err := newID()
	if err != nil {
		return CreateStoryboardOutput{}, fmt.Errorf("failed to create project id")
	}
	jobID, err := newID()
	if err != nil {
		return CreateStoryboardOutput{}, fmt.Errorf("failed to create job id")
	}

	now := time.Now().UTC()
	project := domain.Project{
		ID:                   projectID,
		Title:                strings.TrimSpace(input.Title),
		Style:                strings.TrimSpace(input.Style),
		Platform:             strings.TrimSpace(input.Platform),
		Format:               strings.TrimSpace(input.Format),
		TotalDurationSeconds: input.TotalDurationSeconds,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	job := domain.StoryboardJob{
		ID:          jobID,
		ProjectID:   projectID,
		Status:      "pending",
		CurrentStep: "queued",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assets, err := s.assetStore.StoreProductImages(ctx, projectID, input.ProductImages)
	if err != nil {
		return CreateStoryboardOutput{}, fmt.Errorf("%w: %s", ErrBadInput, err.Error())
	}

	if err := s.repo.SaveProjectBundle(domain.ProjectBundle{
		Project: project,
		Job:     job,
		Assets:  assets,
	}); err != nil {
		return CreateStoryboardOutput{}, fmt.Errorf("save project bundle: %w", err)
	}

	if err := s.queue.EnqueueStoryboardGenerate(ctx, ports.StoryboardGeneratePayload{
		ProjectID: projectID,
		JobID:     jobID,
	}); err != nil {
		_ = s.repo.UpdateJobFailed(projectID, jobID, "queueing", err.Error())
		return CreateStoryboardOutput{}, fmt.Errorf("enqueue storyboard job: %w", err)
	}

	return CreateStoryboardOutput{
		ProjectID: projectID,
		JobID:     jobID,
		Status:    "pending",
	}, nil
}

func (s *StoryboardService) GetStoryboard(projectID string) (domain.ProjectBundle, error) {
	bundle, err := s.repo.GetProjectBundle(projectID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return domain.ProjectBundle{}, ErrNotFound
		}
		return domain.ProjectBundle{}, err
	}

	return bundle, nil
}

func (s *StoryboardService) GetJobStatus(projectID, jobID string) (domain.StoryboardJob, error) {
	job, err := s.repo.GetJob(projectID, jobID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return domain.StoryboardJob{}, ErrNotFound
		}
		return domain.StoryboardJob{}, err
	}

	return job, nil
}
