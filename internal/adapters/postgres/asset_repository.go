package postgres

import "story-board-generator/internal/domain"

func (r *Repository) listAssets(projectID string) []domain.Asset {
	return append([]domain.Asset(nil), r.assets[projectID]...)
}
