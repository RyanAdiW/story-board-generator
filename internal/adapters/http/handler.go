package httpadapter

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"story-board-generator/internal/app"
	"story-board-generator/internal/domain"
)

type Handler struct {
	storyboards *app.StoryboardService
}

func NewHandler(storyboards *app.StoryboardService) *Handler {
	return &Handler{
		storyboards: storyboards,
	}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) CreateStoryboard(c echo.Context) error {
	durationRaw := strings.TrimSpace(c.FormValue("total_duration_seconds"))
	duration, err := strconv.Atoi(durationRaw)
	if err != nil || duration <= 0 {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "total_duration_seconds must be a positive integer",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "multipart form is required",
		})
	}

	files := form.File["product_images[]"]
	if len(files) == 0 {
		files = form.File["product_images"]
	}

	output, err := h.storyboards.CreateStoryboard(c.Request().Context(), app.CreateStoryboardInput{
		Title:                c.FormValue("title"),
		Style:                c.FormValue("style"),
		Platform:             c.FormValue("platform"),
		Format:               c.FormValue("format"),
		TotalDurationSeconds: duration,
		ProductImages:        files,
	})
	if err != nil {
		if errors.Is(err, app.ErrBadInput) {
			return c.JSON(http.StatusBadRequest, errorResponse{
				Message: strings.TrimPrefix(err.Error(), "bad input: "),
			})
		}
		return c.JSON(http.StatusInternalServerError, errorResponse{
			Message: "failed to create storyboard project",
		})
	}

	return c.JSON(http.StatusCreated, createStoryboardResponse{
		ProjectID: output.ProjectID,
		JobID:     output.JobID,
		Status:    output.Status,
	})
}

func (h *Handler) GetStoryboard(c echo.Context) error {
	projectID := c.Param("project_id")

	bundle, err := h.storyboards.GetStoryboard(projectID)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			return c.JSON(http.StatusNotFound, errorResponse{Message: "project not found"})
		}
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to read project"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"project_id":             bundle.Project.ID,
		"title":                  bundle.Project.Title,
		"style":                  bundle.Project.Style,
		"platform":               bundle.Project.Platform,
		"format":                 bundle.Project.Format,
		"total_duration_seconds": bundle.Project.TotalDurationSeconds,
		"assets":                 bundle.Assets,
		"scenes":                 bundle.Scenes,
		"final_image_url":        finalStoryboardURL(bundle.Assets),
	})
}

func (h *Handler) GetJobStatus(c echo.Context) error {
	projectID := c.Param("project_id")
	jobID := c.Param("job_id")

	job, err := h.storyboards.GetJobStatus(projectID, jobID)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			return c.JSON(http.StatusNotFound, errorResponse{Message: "job not found"})
		}
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to read job"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"job_id":        job.ID,
		"project_id":    job.ProjectID,
		"status":        job.Status,
		"current_step":  job.CurrentStep,
		"error_message": job.Error,
	})
}

func finalStoryboardURL(assets []domain.Asset) string {
	for i := len(assets) - 1; i >= 0; i-- {
		if assets[i].AssetType == "final_storyboard_image" {
			return assets[i].FileURL
		}
	}

	return ""
}
