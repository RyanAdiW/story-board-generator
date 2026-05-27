package domain

import "time"

type Asset struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	AssetType string    `json:"asset_type"`
	FileURL   string    `json:"file_url"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}
