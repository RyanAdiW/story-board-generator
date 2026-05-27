package ports

import "context"

type StoryboardGeneratePayload struct {
	ProjectID string `json:"project_id"`
	JobID     string `json:"job_id"`
}

type MessageHandler func(ctx context.Context, payload StoryboardGeneratePayload) error

type QueuePublisher interface {
	EnqueueStoryboardGenerate(ctx context.Context, payload StoryboardGeneratePayload) error
	Close() error
}

type QueueConsumer interface {
	ConsumeStoryboardGenerate(ctx context.Context, handler MessageHandler) error
	Close() error
}
