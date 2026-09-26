package screens

import (
	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
)

func init() { factories["intro"] = (*Scene).introduction }

func (s *Scene) introduction() {
	background, logo, letters := s.image("bg.png"), s.image("logo.png"), s.image("letters.png")
	if s.err != nil {
		return
	}
	stage := s.surface(640, 400)
	text, err := composite.NewHarmonicCellWarp(presets.UnionIntroTextCells())
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
	backdropMotion, err := motion.NewWrapBank(presets.UnionIntroBackdropMotion())
	if err != nil {
		s.err = err
		return
	}
	logoMotion, err := motion.NewHarmonicTransform(presets.UnionIntroLogoTransform())
	if err != nil {
		s.err = err
		return
	}
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		// Repetition keeps the background moving without an exposed seam.
		back.DrawAt(stage, background, 0, backdropMotion.At(0))
		backdropMotion.Step()
		pose := logoMotion.At(float64(s.Frame) * presets.UnionIntroLogoPhaseStep())
		logoOptions := ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		logoOptions.GeoM.Scale(pose.ScaleX, pose.ScaleY)
		logoOptions.GeoM.Translate(pose.X, pose.Y)
		composite.Instance{Image: logo, Options: logoOptions}.Draw(stage)
		text.DrawAt(stage, letters, kit.Frame{Tick: s.Frame}, 70, 136)
		s.draw(s.Canvas, stage, 64, 68)
	}
}
