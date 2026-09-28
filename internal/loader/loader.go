// Package loader presents the animated credits between the hall and a screen.
package loader

import (
	"fmt"
	"image/color"
	"io/fs"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/assets"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/timeline"
)

// Duration follows the embedded transition recording.
const Duration = 3777120 * time.Microsecond

type Screen struct {
	store       *assets.Store
	font, paper *ebiten.Image
	reveal      *scrolling.Scrolling
	clock       *timeline.CueClock
}

func New(id string, files fs.FS, rate int) (*Screen, error) {
	lines, ok := credits[id]
	if !ok {
		return nil, fmt.Errorf("loader: unknown screen %q", id)
	}
	clock, err := timeline.NewCueClock(timeline.CueClockConfig{Rate: rate})
	if err != nil {
		return nil, err
	}
	s := &Screen{store: assets.New(files), clock: clock}
	s.font, err = s.store.Texture("loader/loader.png")
	if err != nil {
		s.store.Close()
		return nil, err
	}
	grid, err := presets.BitmapFont("union-loader", s.font, ebiten.FilterNearest)
	if err != nil {
		s.store.Close()
		return nil, err
	}
	config := presets.UnionCreditsReveal(grid, lines)
	s.reveal, err = scrolling.New(scrolling.Config{Reveal: &scrolling.RevealTransportConfig{
		Reveal: config,
		TimeAt: func(frame kit.Frame) float64 { return float64(frame.Tick) * 140 * 60 / float64(clock.Rate()) },
	}})
	if err != nil {
		s.store.Close()
		return nil, err
	}
	s.paper = ebiten.NewImage(640, 400)
	return s, nil
}
func (s *Screen) Update() {
	s.clock.Step()
	_ = s.reveal.Update(kit.Frame{Tick: uint64(s.clock.Tick())})
}
func (s *Screen) Done() bool {
	return s.clock.Reached(Duration)
}
func (s *Screen) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	s.paper.Clear()
	// The shared reveal uses a fixed clock independent from display refresh.
	s.reveal.Draw(s.paper)

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(64, 60)
	dst.DrawImage(s.paper, &op)
}
func (s *Screen) Close() error {
	s.reveal.Close()
	s.paper.Deallocate()
	return s.store.Close()
}
