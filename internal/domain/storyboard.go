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

type ProjectBundle struct {
	Project Project       `json:"project"`
	Job     StoryboardJob `json:"job"`
	Assets  []Asset       `json:"assets"`
}
