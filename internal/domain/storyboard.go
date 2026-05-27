package domain

import "time"

type Project struct {
	ID                   string    `json:"id"`
	Title                string    `json:"title"`
	Style                string    `json:"style"`
	Platform             string    `json:"platform"`
	Format               string    `json:"format"`
	TotalDurationSeconds int       `json:"total_duration_seconds"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type Asset struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	AssetType string    `json:"asset_type"`
	FileURL   string    `json:"file_url"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

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

type ProjectBundle struct {
	Project Project       `json:"project"`
	Job     StoryboardJob `json:"job"`
	Assets  []Asset       `json:"assets"`
}
