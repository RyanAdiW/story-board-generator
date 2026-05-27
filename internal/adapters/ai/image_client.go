package ai

import (
	"context"
	"fmt"
)

type ImageClient struct{}

func NewImageClient() *ImageClient {
	return &ImageClient{}
}

func (c *ImageClient) Generate(_ context.Context, _ string) (string, error) {
	return "", fmt.Errorf("image client is not implemented yet")
}
