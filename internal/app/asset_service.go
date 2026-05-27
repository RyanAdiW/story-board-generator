package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"story-board-generator/internal/domain"
	"story-board-generator/internal/ports"
)

var allowedExtensions = map[string]struct{}{
	".png":  {},
	".jpg":  {},
	".jpeg": {},
	".webp": {},
}

const maxImages = 10

type AssetService struct {
	storage ports.Storage
}

func NewAssetService(storage ports.Storage) *AssetService {
	return &AssetService{storage: storage}
}

func (s *AssetService) StoreProductImages(ctx context.Context, projectID string, files []*multipart.FileHeader) ([]domain.Asset, error) {
	if len(files) == 0 || len(files) > maxImages {
		return nil, fmt.Errorf("product_images must contain between 1 and 10 files")
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

		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("open file %s: %w", fileHeader.Filename, err)
		}

		assetID, err := newID()
		if err != nil {
			_ = file.Close()
			return nil, fmt.Errorf("failed to generate asset id")
		}

		objectPath := filepath.ToSlash(filepath.Join(projectID, "inputs", assetID+ext))
		fileURL, err := s.storage.Upload(ctx, objectPath, contentType, file)
		_ = file.Close()
		if err != nil {
			return nil, fmt.Errorf("store %s: %w", fileHeader.Filename, err)
		}

		assets = append(assets, domain.Asset{
			ID:        assetID,
			ProjectID: projectID,
			AssetType: "input_product_image",
			FileURL:   fileURL,
			MimeType:  contentType,
			CreatedAt: now,
		})
	}

	return assets, nil
}

func newID() (string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	return hex.EncodeToString(raw), nil
}
