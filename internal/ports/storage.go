package ports

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, objectPath, contentType string, reader io.Reader) (string, error)
}
