package renderer

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"sort"
	"strings"

	"story-board-generator/internal/domain"
	"story-board-generator/internal/ports"
)

type StoryboardRenderer struct{}

func NewStoryboardRenderer() *StoryboardRenderer {
	return &StoryboardRenderer{}
}

func (r *StoryboardRenderer) RenderStoryboard(_ context.Context, input ports.RenderStoryboardInput) (ports.RenderStoryboardOutput, error) {
	if len(input.Scenes) == 0 {
		return ports.RenderStoryboardOutput{}, fmt.Errorf("no scenes to render")
	}

	scenes := append([]struct {
		SceneNumber int
		ImageURL    string
	}{}, collectSceneRefs(input.Scenes)...)
	sort.Slice(scenes, func(i, j int) bool { return scenes[i].SceneNumber < scenes[j].SceneNumber })

	const (
		canvasW = 2048
		margin  = 48
		gap     = 28
		headerH = 140
		cols    = 2
	)

	cellW := (canvasW - margin*2 - gap*(cols-1)) / cols
	cellH := int(float64(cellW) * 9.0 / 16.0)
	rows := (len(scenes) + cols - 1) / cols
	canvasH := margin + headerH + gap + rows*cellH + (rows-1)*gap + margin

	canvas := image.NewRGBA(image.Rect(0, 0, canvasW, canvasH))
	bg := color.RGBA{R: 18, G: 20, B: 24, A: 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)

	headerRect := image.Rect(margin, margin, canvasW-margin, margin+headerH)
	draw.Draw(canvas, headerRect, &image.Uniform{C: color.RGBA{R: 34, G: 38, B: 46, A: 255}}, image.Point{}, draw.Src)

	for idx, ref := range scenes {
		col := idx % cols
		row := idx / cols
		x := margin + col*(cellW+gap)
		y := margin + headerH + gap + row*(cellH+gap)
		cell := image.Rect(x, y, x+cellW, y+cellH)

		draw.Draw(canvas, cell, &image.Uniform{C: color.RGBA{R: 40, G: 44, B: 54, A: 255}}, image.Point{}, draw.Src)

		img, err := loadImage(ref.ImageURL)
		if err != nil || img == nil {
			draw.Draw(canvas, inset(cell, 8), &image.Uniform{C: color.RGBA{R: 58, G: 62, B: 72, A: 255}}, image.Point{}, draw.Src)
			continue
		}
		drawImageCover(canvas, inset(cell, 8), img)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, canvas); err != nil {
		return ports.RenderStoryboardOutput{}, fmt.Errorf("encode storyboard image: %w", err)
	}

	return ports.RenderStoryboardOutput{
		Bytes:    buf.Bytes(),
		MimeType: "image/png",
	}, nil
}

func collectSceneRefs(scenes []domain.Scene) []struct {
	SceneNumber int
	ImageURL    string
} {
	out := make([]struct {
		SceneNumber int
		ImageURL    string
	}, 0, len(scenes))
	for _, scene := range scenes {
		out = append(out, struct {
			SceneNumber int
			ImageURL    string
		}{
			SceneNumber: scene.SceneNumber,
			ImageURL:    scene.ImageURL,
		})
	}
	return out
}

func loadImage(path string) (image.Image, error) {
	filePath := strings.TrimSpace(path)
	if filePath == "" {
		return nil, fmt.Errorf("empty image path")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if strings.HasSuffix(strings.ToLower(filePath), ".png") {
		return png.Decode(f)
	}
	if strings.HasSuffix(strings.ToLower(filePath), ".jpg") || strings.HasSuffix(strings.ToLower(filePath), ".jpeg") {
		return jpeg.Decode(f)
	}

	img, _, err := image.Decode(f)
	return img, err
}

func drawImageCover(dst draw.Image, target image.Rectangle, src image.Image) {
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	tw, th := target.Dx(), target.Dy()
	if sw == 0 || sh == 0 || tw <= 0 || th <= 0 {
		return
	}

	scale := max(float64(tw)/float64(sw), float64(th)/float64(sh))
	cw := int(float64(tw) / scale)
	ch := int(float64(th) / scale)
	if cw <= 0 {
		cw = 1
	}
	if ch <= 0 {
		ch = 1
	}

	srcX0 := sb.Min.X + (sw-cw)/2
	srcY0 := sb.Min.Y + (sh-ch)/2

	for y := 0; y < th; y++ {
		sy := srcY0 + (y*ch)/th
		for x := 0; x < tw; x++ {
			sx := srcX0 + (x*cw)/tw
			dst.Set(target.Min.X+x, target.Min.Y+y, src.At(sx, sy))
		}
	}
}

func inset(r image.Rectangle, pad int) image.Rectangle {
	return image.Rect(r.Min.X+pad, r.Min.Y+pad, r.Max.X-pad, r.Max.Y-pad)
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
