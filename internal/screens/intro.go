package screens

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
)

func init() { factories["intro"] = (*Scene).introduction }

func (s *Scene) introduction() {
	background, logo, letters := s.image("bg.png"), s.image("logo.png"), s.image("letters.png")
	if s.err != nil {
		return
	}
	stage := s.surface(640, 400)
	text, err := composite.NewHarmonicCellWarp(composite.HarmonicCellWarpConfig{
		Cell: image.Pt(32, 16), Filter: ebiten.FilterLinear,
		Waves: composite.CellWaveBank{
			XRows:    motion.Waves{{Amplitude: 32, Spatial: .3, Speed: .08}},
			YColumns: motion.Waves{{Amplitude: 16, Spatial: .3, Speed: .08}},
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
