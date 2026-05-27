package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"story-board-generator/internal/domain"
)

const metadataFilename = "metadata.json"

var ErrProjectNotFound = errors.New("project not found")

type Store struct {
	path string

	mu       sync.RWMutex
	projects map[string]domain.Project
	jobs     map[string]domain.StoryboardJob
	assets   map[string][]domain.Asset
}

type snapshot struct {
	Projects map[string]domain.Project       `json:"projects"`
	Jobs     map[string]domain.StoryboardJob `json:"jobs"`
	Assets   map[string][]domain.Asset       `json:"assets"`
}

func New(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	store := &Store{
		path:     filepath.Join(dataDir, metadataFilename),
		projects: map[string]domain.Project{},
		jobs:     map[string]domain.StoryboardJob{},
		assets:   map[string][]domain.Asset{},
	}

	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) SaveProjectBundle(bundle domain.ProjectBundle) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.projects[bundle.Project.ID] = bundle.Project
	s.jobs[bundle.Job.ID] = bundle.Job
	s.assets[bundle.Project.ID] = append([]domain.Asset(nil), bundle.Assets...)

	return s.persist()
}

func (s *Store) GetProjectBundle(projectID string) (domain.ProjectBundle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, ok := s.projects[projectID]
	if !ok {
		return domain.ProjectBundle{}, ErrProjectNotFound
	}

	var projectJob domain.StoryboardJob
	for _, job := range s.jobs {
		if job.ProjectID == projectID {
			projectJob = job
			break
		}
	}

	return domain.ProjectBundle{
		Project: project,
		Job:     projectJob,
		Assets:  append([]domain.Asset(nil), s.assets[projectID]...),
	}, nil
}

func (s *Store) GetJob(projectID, jobID string) (domain.StoryboardJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[jobID]
	if !ok || job.ProjectID != projectID {
		return domain.StoryboardJob{}, ErrProjectNotFound
	}

	return job, nil
}

func (s *Store) load() error {
	content, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read metadata: %w", err)
	}

	var snap snapshot
	if err := json.Unmarshal(content, &snap); err != nil {
		return fmt.Errorf("parse metadata: %w", err)
	}

	if snap.Projects != nil {
		s.projects = snap.Projects
	}
	if snap.Jobs != nil {
		s.jobs = snap.Jobs
	}
	if snap.Assets != nil {
		s.assets = snap.Assets
	}

	return nil
}

func (s *Store) persist() error {
	snap := snapshot{
		Projects: s.projects,
		Jobs:     s.jobs,
		Assets:   s.assets,
	}

	encoded, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	if err := os.WriteFile(s.path, encoded, 0o644); err != nil {
		return fmt.Errorf("write metadata: %w", err)
	}

	return nil
}
