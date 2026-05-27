package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type S3Storage struct {
	baseDir string
}

func NewS3Storage(baseDir string) *S3Storage {
	return &S3Storage{baseDir: baseDir}
}

func (s *S3Storage) Upload(_ context.Context, objectPath, _ string, reader io.Reader) (string, error) {
	target := filepath.Join(s.baseDir, filepath.FromSlash(objectPath))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", fmt.Errorf("create storage directory: %w", err)
	}

	file, err := os.Create(target)
	if err != nil {
		return "", fmt.Errorf("create target file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("write target file: %w", err)
	}

	return filepath.ToSlash(target), nil
}
