package ai

import (
	"context"
	"fmt"
)

type OpenAIClient struct{}

func NewOpenAIClient() *OpenAIClient {
	return &OpenAIClient{}
}

func (c *OpenAIClient) GenerateStoryboard(_ context.Context, _ string) (string, error) {
	return "", fmt.Errorf("openai client is not implemented yet")
}
