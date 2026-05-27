package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"strings"
	"time"

	"story-board-generator/internal/ports"
)

type ImageClient struct {
	apiKey     string
	model      string
	size       string
	quality    string
	httpClient *http.Client
}

func NewImageClient(apiKey, model, size, quality string) *ImageClient {
	if strings.TrimSpace(model) == "" {
		model = "gpt-image-1"
	}
	if strings.TrimSpace(size) == "" {
		size = "1024x1536"
	}
	if strings.TrimSpace(quality) == "" {
		quality = "medium"
	}

	return &ImageClient{
		apiKey:  strings.TrimSpace(apiKey),
		model:   model,
		size:    size,
		quality: quality,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *ImageClient) GenerateSceneImage(ctx context.Context, input ports.SceneImageInput) (ports.SceneImageOutput, error) {
	if c.apiKey == "" {
		return fallbackImage(), nil
	}

	prompt := input.Scene.ImagePrompt
	if strings.TrimSpace(prompt) == "" {
		prompt = fmt.Sprintf(
			"%s. Style: %s. Format: %s.",
			input.Scene.VisualDescription,
			input.Project.Style,
			input.Project.Format,
		)
	}

	reqBody := map[string]any{
		"model":   c.model,
		"prompt":  prompt,
		"size":    c.size,
		"quality": c.quality,
	}

	raw, err := json.Marshal(reqBody)
	if err != nil {
		return ports.SceneImageOutput{}, fmt.Errorf("marshal image request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/images/generations", bytes.NewReader(raw))
	if err != nil {
		return ports.SceneImageOutput{}, fmt.Errorf("build image request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ports.SceneImageOutput{}, fmt.Errorf("send image request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ports.SceneImageOutput{}, fmt.Errorf("read image response: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return ports.SceneImageOutput{}, fmt.Errorf("image api error: %s", string(respBytes))
	}

	var parsed struct {
		Data []struct {
			B64JSON string `json:"b64_json"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return ports.SceneImageOutput{}, fmt.Errorf("decode image response: %w", err)
	}
	if len(parsed.Data) == 0 || parsed.Data[0].B64JSON == "" {
		return ports.SceneImageOutput{}, fmt.Errorf("empty image response data")
	}

	imgBytes, err := base64.StdEncoding.DecodeString(parsed.Data[0].B64JSON)
	if err != nil {
		return ports.SceneImageOutput{}, fmt.Errorf("decode base64 image: %w", err)
	}

	return ports.SceneImageOutput{
		Bytes:    imgBytes,
		MimeType: "image/png",
	}, nil
}

func fallbackImage() ports.SceneImageOutput {
	img := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	bg := color.RGBA{R: 30, G: 30, B: 35, A: 255}
	for y := 0; y < 1024; y++ {
		for x := 0; x < 1024; x++ {
			img.Set(x, y, bg)
		}
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)

	return ports.SceneImageOutput{
		Bytes:    buf.Bytes(),
		MimeType: "image/png",
	}
}
