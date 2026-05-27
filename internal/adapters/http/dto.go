package httpadapter

type createStoryboardResponse struct {
	ProjectID string `json:"project_id"`
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
}

type errorResponse struct {
	Message string `json:"message"`
}
