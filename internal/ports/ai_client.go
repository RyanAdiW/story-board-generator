package ports

import "context"

type AIClient interface {
	GenerateStoryboard(ctx context.Context, prompt string) (string, error)
}
