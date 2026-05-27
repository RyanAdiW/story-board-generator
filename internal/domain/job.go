package domain

import "time"

type StoryboardJob struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Status      string     `json:"status"`
	CurrentStep string     `json:"current_step,omitempty"`
	Error       string     `json:"error_message,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
