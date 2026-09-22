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
	text, err := composite.NewCellWarp(composite.CellWarpConfig{
		Cell: image.Pt(32, 16), Filter: ebiten.FilterLinear,
		Sample: func(row, column int, frame kit.Frame) composite.CellTransform {
			phase := float64(frame.Tick) * .08
			pose := composite.CellTransform{}
			pose.GeoM.Translate(32*math.Sin(phase+float64(row)*.3), 16*math.Sin(phase+float64(column)*.3))
			return pose
		},
	})
	if err != nil {
		s.err = err
		return
	}
	s.closers = append(s.closers, text.Close)
	back, err := composite.NewBackground(composite.BackgroundConfig{PeriodX: 640, PeriodY: 96})
	if err != nil {
		s.err = err
		return
	}
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		// Repetition keeps the background moving without an exposed seam.
		back.DrawAt(stage, background, 0, float64((s.Frame%24)*4))
		logoOptions := ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		logoOptions.GeoM.Translate(256+190*math.Sin(float64(s.Frame)*.05), 43)
		composite.Instance{Image: logo, Options: logoOptions}.Draw(stage)
		text.DrawAt(stage, letters, kit.Frame{Tick: s.Frame}, 70, 136)
		s.draw(s.Canvas, stage, 64, 68)
	}
}
