package ports

import (
	"errors"

	"story-board-generator/internal/domain"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	SaveProjectBundle(bundle domain.ProjectBundle) error
	GetProjectBundle(projectID string) (domain.ProjectBundle, error)
	GetJob(projectID, jobID string) (domain.StoryboardJob, error)
	UpdateJobProcessing(projectID, jobID, step string) error
	UpdateJobCompleted(projectID, jobID string) error
	UpdateJobFailed(projectID, jobID, step, errMessage string) error
}
