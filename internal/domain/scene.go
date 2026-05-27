package domain

type Scene struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	SceneNumber int    `json:"scene_number"`
}
