package screens

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
)

func init() { factories["intro"] = (*Scene).introduction }

func (s *Scene) introduction() {
	background, logo, letters := s.image("bg.png"), s.image("logo.png"), s.image("letters.png")
	if s.err != nil {
		return
	}
	stage := s.surface(640, 400)
	text := introTiles(letters, image.Pt(32, 16), image.Pt(70, 136))
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		// Repetition keeps the background moving without an exposed seam.
		composite.Repeat(stage, background, composite.Repetition{Zoom: 1, PhaseY: -float64((s.Frame % 24) * 4)})
		logoOptions := ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		logoOptions.GeoM.Translate(256+190*math.Sin(float64(s.Frame)*.05), 43)
		composite.Instance{Image: logo, Options: logoOptions}.Draw(stage)
		text.Update(kit.Frame{Tick: s.Frame})
		text.Draw(stage)
		s.draw(s.Canvas, stage, 64, 68)
	}
}

// introTiles applies independent row and column waves to image fragments. Each
// fragment keeps its source column when a row shifts, preserving the crossed
// motion. DCK's sprite sampler works with any image and cell dimensions.
func introTiles(source *ebiten.Image, cell, origin image.Point) *composite.Sprites {
	bounds := source.Bounds()
	columns := (bounds.Dx() + cell.X - 1) / cell.X
	rows := (bounds.Dy() + cell.Y - 1) / cell.Y
	tiles := make([]image.Rectangle, columns*rows)
	for i := range tiles {
		left, top := i%columns*cell.X, i/columns*cell.Y
		tiles[i] = image.Rect(left, top, min(left+cell.X, bounds.Dx()), min(top+cell.Y, bounds.Dy())).Add(bounds.Min)
	}
	return &composite.Sprites{
		Count: len(tiles),
		Sample: func(index int, frame kit.Frame) composite.Instance {
			column, row := index%columns, index/columns
			phase := float64(frame.Tick) * .08
			x := float64(origin.X+column*cell.X) + 32*math.Sin(phase+float64(row)*.3)
			y := float64(origin.Y+row*cell.Y) + 16*math.Sin(phase+float64(column)*.3)
			op := ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
			op.GeoM.Translate(x, y)
			return composite.Instance{Image: source, Source: &tiles[index], Options: op}
		},
	}
}
