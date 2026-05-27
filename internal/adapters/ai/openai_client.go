package ai

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"story-board-generator/internal/domain"
	"story-board-generator/internal/ports"
)

type OpenAIClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	if strings.TrimSpace(model) == "" {
		model = "gpt-4.1-mini"
	}

	return &OpenAIClient{
		apiKey: strings.TrimSpace(apiKey),
		model:  model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *OpenAIClient) GenerateScenes(ctx context.Context, input ports.SceneGenerationInput) ([]domain.Scene, error) {
	if c.apiKey == "" {
		return fallbackScenes(input.Project), nil
	}

	reqBody := map[string]any{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert ad storyboard director. Return valid JSON only.",
			},
			{
				"role":    "user",
				"content": buildScenePrompt(input.Project),
			},
		},
		"response_format": map[string]string{
			"type": "json_object",
		},
		"temperature": 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("openai api error: %s", string(respBytes))
	}

	var completion chatCompletionResponse
	if err := json.Unmarshal(respBytes, &completion); err != nil {
		return nil, fmt.Errorf("decode completion: %w", err)
	}
	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("empty completion choices")
	}

	content := completion.Choices[0].Message.Content
	var out sceneOutput
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return nil, fmt.Errorf("decode scene json: %w", err)
	}
	if len(out.Scenes) == 0 {
		return nil, fmt.Errorf("no scenes returned")
	}

	scenes := make([]domain.Scene, 0, len(out.Scenes))
	now := time.Now().UTC()
	for i, raw := range out.Scenes {
		scenes = append(scenes, domain.Scene{
			ID:                mustID(),
			ProjectID:         input.Project.ID,
			SceneNumber:       i + 1,
			StartSecond:       raw.StartSecond,
			EndSecond:         raw.EndSecond,
			VisualDescription: strings.TrimSpace(raw.VisualDescription),
			CameraDirection:   strings.TrimSpace(raw.CameraDirection),
			MotionDescription: strings.TrimSpace(raw.MotionDescription),
			SoundFX:           strings.TrimSpace(raw.SoundFX),
			OnScreenText:      strings.TrimSpace(raw.OnScreenText),
			Notes:             strings.TrimSpace(raw.Notes),
			ImagePrompt:       strings.TrimSpace(raw.ImagePrompt),
			CreatedAt:         now,
		})
	}

	return scenes, nil
}

type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type sceneOutput struct {
	Scenes []scenePayload `json:"scenes"`
}

type scenePayload struct {
	StartSecond       int    `json:"start_second"`
	EndSecond         int    `json:"end_second"`
	VisualDescription string `json:"visual_description"`
	CameraDirection   string `json:"camera_direction"`
	MotionDescription string `json:"motion_description"`
	SoundFX           string `json:"sound_fx"`
	OnScreenText      string `json:"on_screen_text"`
	Notes             string `json:"notes"`
	ImagePrompt       string `json:"image_prompt"`
}

func buildScenePrompt(project domain.Project) string {
	return fmt.Sprintf(
		"Generate cinematic storyboard scenes as JSON for this ad. title=%q style=%q platform=%q format=%q duration_seconds=%d. Return object: {\"scenes\":[...]} with fields: start_second,end_second,visual_description,camera_direction,motion_description,sound_fx,on_screen_text,notes,image_prompt. Keep 3-8 scenes and cover entire duration.",
		project.Title,
		project.Style,
		project.Platform,
		project.Format,
		project.TotalDurationSeconds,
	)
}

func fallbackScenes(project domain.Project) []domain.Scene {
	sceneCount := project.TotalDurationSeconds / 5
	if sceneCount < 3 {
		sceneCount = 3
	}
	if sceneCount > 8 {
		sceneCount = 8
	}

	step := project.TotalDurationSeconds / sceneCount
	if step <= 0 {
		step = 1
	}

	now := time.Now().UTC()
	scenes := make([]domain.Scene, 0, sceneCount)
	for i := 0; i < sceneCount; i++ {
		start := i * step
		end := start + step
		if i == sceneCount-1 {
			end = project.TotalDurationSeconds
		}
		if end <= start {
			end = start + 1
		}

		scenes = append(scenes, domain.Scene{
			ID:                mustID(),
			ProjectID:         project.ID,
			SceneNumber:       i + 1,
			StartSecond:       start,
			EndSecond:         end,
			VisualDescription: fmt.Sprintf("%s product hero shot in %s style.", project.Title, project.Style),
			CameraDirection:   "Cinematic close-up with smooth push-in",
			MotionDescription: "Subtle movement with dramatic reveal pacing",
			SoundFX:           "Cinematic whoosh and bass impact",
			OnScreenText:      fmt.Sprintf("SCENE %d", i+1),
			Notes:             fmt.Sprintf("Optimized for %s in %s format.", project.Platform, project.Format),
			ImagePrompt:       fmt.Sprintf("Use product references, %s style, %s aspect, high-detail advertising frame, no watermark.", project.Style, project.Format),
			CreatedAt:         now,
		})
	}

	return scenes
}

func mustID() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(raw)
}
