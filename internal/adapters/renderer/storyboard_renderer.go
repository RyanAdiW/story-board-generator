package renderer

import "fmt"

type StoryboardRenderer struct{}

func NewStoryboardRenderer() *StoryboardRenderer {
	return &StoryboardRenderer{}
}

func (r *StoryboardRenderer) Render() error {
	return fmt.Errorf("storyboard renderer is not implemented yet")
}
