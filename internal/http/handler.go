package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"story-board-generator/internal/domain"
	"story-board-generator/internal/store"
)

var allowedExtensions = map[string]struct{}{
	".png":  {},
	".jpg":  {},
	".jpeg": {},
	".webp": {},
}

const maxImages = 10

type Handler struct {
	repo       *store.Store
	uploadsDir string
}

func NewHandler(repo *store.Store, uploadsDir string) *Handler {
	return &Handler{
		repo:       repo,
		uploadsDir: uploadsDir,
	}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) CreateStoryboard(c echo.Context) error {
	title := strings.TrimSpace(c.FormValue("title"))
	style := strings.TrimSpace(c.FormValue("style"))
	platform := strings.TrimSpace(c.FormValue("platform"))
	format := strings.TrimSpace(c.FormValue("format"))
	durationRaw := strings.TrimSpace(c.FormValue("total_duration_seconds"))

	if title == "" || style == "" || platform == "" || format == "" || durationRaw == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "title, style, platform, format, and total_duration_seconds are required",
		})
	}

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

	if len(files) == 0 || len(files) > maxImages {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "product_images must contain between 1 and 10 files",
		})
	}

	projectID, err := newID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to create project id"})
	}
	jobID, err := newID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to create job id"})
	}

	now := time.Now().UTC()
	project := domain.Project{
		ID:                   projectID,
		Title:                title,
		Style:                style,
		Platform:             platform,
		Format:               format,
		TotalDurationSeconds: duration,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	job := domain.StoryboardJob{
		ID:          jobID,
		ProjectID:   projectID,
		Status:      "pending",
		CurrentStep: "uploading_assets",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assets, err := h.saveProductImages(projectID, files)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: err.Error()})
	}

	if err := h.repo.SaveProjectBundle(domain.ProjectBundle{
		Project: project,
		Job:     job,
		Assets:  assets,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{
			Message: "failed to persist project metadata",
		})
	}

	return c.JSON(http.StatusCreated, createStoryboardResponse{
		ProjectID: projectID,
		JobID:     jobID,
		Status:    "pending",
	})
}

func (h *Handler) GetStoryboard(c echo.Context) error {
	projectID := c.Param("project_id")

	bundle, err := h.repo.GetProjectBundle(projectID)
	if err != nil {
		if err == store.ErrProjectNotFound {
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
		"final_image_url":        "",
	})
}

func (h *Handler) GetJobStatus(c echo.Context) error {
	projectID := c.Param("project_id")
	jobID := c.Param("job_id")

	job, err := h.repo.GetJob(projectID, jobID)
	if err != nil {
		if err == store.ErrProjectNotFound {
			return c.JSON(http.StatusNotFound, errorResponse{Message: "job not found"})
		}
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to read job"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"job_id":        job.ID,
		"project_id":    job.ProjectID,
		"status":        job.Status,
		"current_step":  job.CurrentStep,
		"error_message": null,
	})
}

func (h *Handler) saveProductImages(projectID string, files []*multipart.FileHeader) ([]domain.Asset, error) {
	projectDir := filepath.Join(h.uploadsDir, projectID, "inputs")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory")
	}

	now := time.Now().UTC()
	assets := make([]domain.Asset, 0, len(files))

	for _, fileHeader := range files {
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if _, ok := allowedExtensions[ext]; !ok {
			return nil, fmt.Errorf("unsupported image extension: %s", ext)
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			return nil, fmt.Errorf("file %s must be an image", fileHeader.Filename)
		}

		assetID, err := newID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate asset id")
		}

		filename := assetID + ext
		target := filepath.Join(projectDir, filename)
		if err := copyFile(fileHeader, target); err != nil {
			return nil, fmt.Errorf("failed to store %s", fileHeader.Filename)
		}

		assets = append(assets, domain.Asset{
			ID:        assetID,
			ProjectID: projectID,
			AssetType: "input_product_image",
			FileURL:   filepath.ToSlash(target),
			MimeType:  contentType,
			CreatedAt: now,
		})
	}

	return assets, nil
}

func copyFile(fileHeader *multipart.FileHeader, target string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func newID() (string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}
