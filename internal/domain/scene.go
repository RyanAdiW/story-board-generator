package domain

import "time"

type Scene struct {
	ID                string    `json:"id"`
	ProjectID         string    `json:"project_id"`
	SceneNumber       int       `json:"scene_number"`
	StartSecond       int       `json:"start_second"`
	EndSecond         int       `json:"end_second"`
	VisualDescription string    `json:"visual_description"`
	CameraDirection   string    `json:"camera_direction"`
	MotionDescription string    `json:"motion_description"`
	SoundFX           string    `json:"sound_fx"`
	OnScreenText      string    `json:"on_screen_text"`
	Notes             string    `json:"notes,omitempty"`
	ImagePrompt       string    `json:"image_prompt"`
	ImageAssetID      string    `json:"image_asset_id,omitempty"`
	ImageURL          string    `json:"image_url,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}
