package postgres

import (
	"time"

	"story-board-generator/internal/domain"
	"story-board-generator/internal/ports"
)

func (r *Repository) GetJob(projectID, jobID string) (domain.StoryboardJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, ok := r.jobs[jobID]
	if !ok || job.ProjectID != projectID {
		return domain.StoryboardJob{}, errNotFound()
	}

	return job, nil
}

func (r *Repository) UpdateJobProcessing(projectID, jobID, step string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, ok := r.jobs[jobID]
	if !ok || job.ProjectID != projectID {
		return errNotFound()
	}

	job.Status = "processing"
	job.CurrentStep = step
	job.Error = ""
	job.UpdatedAt = nowUTC()
	job.CompletedAt = nil
	r.jobs[jobID] = job

	return r.persist()
}

func (r *Repository) UpdateJobCompleted(projectID, jobID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, ok := r.jobs[jobID]
	if !ok || job.ProjectID != projectID {
		return errNotFound()
	}

	completedAt := nowUTC()
	job.Status = "completed"
	job.CurrentStep = "completed"
	job.Error = ""
	job.UpdatedAt = completedAt
	job.CompletedAt = &completedAt
	r.jobs[jobID] = job

	return r.persist()
}

func (r *Repository) UpdateJobFailed(projectID, jobID, step, errMessage string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, ok := r.jobs[jobID]
	if !ok || job.ProjectID != projectID {
		return errNotFound()
	}

	job.Status = "failed"
	job.CurrentStep = step
	job.Error = errMessage
	job.UpdatedAt = nowUTC()
	job.CompletedAt = nil
	r.jobs[jobID] = job

	return r.persist()
}

func nowUTC() time.Time {
	return time.Now().UTC()
}

func errNotFound() error {
	return ports.ErrNotFound
}
