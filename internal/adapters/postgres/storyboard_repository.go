package postgres

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"story-board-generator/internal/domain"
)

const metadataFilename = "metadata.json"

type Repository struct {
	path string

	mu       sync.RWMutex
	projects map[string]domain.Project
	jobs     map[string]domain.StoryboardJob
	assets   map[string][]domain.Asset
	scenes   map[string][]domain.Scene
}

type snapshot struct {
	Projects map[string]domain.Project       `json:"projects"`
	Jobs     map[string]domain.StoryboardJob `json:"jobs"`
	Assets   map[string][]domain.Asset       `json:"assets"`
	Scenes   map[string][]domain.Scene       `json:"scenes"`
}

func NewRepository(dataDir string) (*Repository, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	repo := &Repository{
		path:     filepath.Join(dataDir, metadataFilename),
		projects: map[string]domain.Project{},
		jobs:     map[string]domain.StoryboardJob{},
		assets:   map[string][]domain.Asset{},
		scenes:   map[string][]domain.Scene{},
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) SaveProjectBundle(bundle domain.ProjectBundle) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.projects[bundle.Project.ID] = bundle.Project
	r.jobs[bundle.Job.ID] = bundle.Job
	r.assets[bundle.Project.ID] = append([]domain.Asset(nil), bundle.Assets...)
	if bundle.Scenes != nil {
		r.scenes[bundle.Project.ID] = append([]domain.Scene(nil), bundle.Scenes...)
	}

	return r.persist()
}

func (r *Repository) SaveScenes(projectID string, scenes []domain.Scene) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.projects[projectID]; !ok {
		return errNotFound()
	}

	r.scenes[projectID] = append([]domain.Scene(nil), scenes...)
	return r.persist()
}

func (r *Repository) GetProjectBundle(projectID string) (domain.ProjectBundle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	project, ok := r.projects[projectID]
	if !ok {
		return domain.ProjectBundle{}, errNotFound()
	}

	var projectJob domain.StoryboardJob
	for _, job := range r.jobs {
		if job.ProjectID == projectID {
			projectJob = job
			break
		}
	}

	return domain.ProjectBundle{
		Project: project,
		Job:     projectJob,
		Assets:  append([]domain.Asset(nil), r.assets[projectID]...),
		Scenes:  append([]domain.Scene(nil), r.scenes[projectID]...),
	}, nil
}

func (r *Repository) load() error {
	content, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read metadata: %w", err)
	}

	var snap snapshot
	if err := json.Unmarshal(content, &snap); err != nil {
		return fmt.Errorf("parse metadata: %w", err)
	}

	if snap.Projects != nil {
		r.projects = snap.Projects
	}
	if snap.Jobs != nil {
		r.jobs = snap.Jobs
	}
	if snap.Assets != nil {
		r.assets = snap.Assets
	}
	if snap.Scenes != nil {
		r.scenes = snap.Scenes
	}

	return nil
}

func (r *Repository) persist() error {
	snap := snapshot{
		Projects: r.projects,
		Jobs:     r.jobs,
		Assets:   r.assets,
		Scenes:   r.scenes,
	}

	encoded, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	if err := os.WriteFile(r.path, encoded, 0o644); err != nil {
		return fmt.Errorf("write metadata: %w", err)
	}

	return nil
}
