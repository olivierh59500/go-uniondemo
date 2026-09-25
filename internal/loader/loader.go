// Package loader presents the animated credits between the hall and a screen.
package loader

import (
	"fmt"
	"image/color"
	"io/fs"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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
	reveal      *scrolling.Reveal
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
	s.reveal, err = scrolling.NewReveal(config)
	if err != nil {
		s.store.Close()
		return nil, err
	}
	s.paper = ebiten.NewImage(640, 400)
	return s, nil
}
func (s *Screen) Update() { s.clock.Step() }
func (s *Screen) Done() bool {
	return s.clock.Reached(Duration)
}
func (s *Screen) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	s.paper.Clear()
	// The shared reveal uses a fixed clock independent from display refresh.
	s.reveal.DrawAt(s.paper, float64(s.clock.Tick())*140*60/float64(s.clock.Rate()))

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(64, 60)
	dst.DrawImage(s.paper, &op)
}
func (s *Screen) Close() error { s.paper.Deallocate(); return s.store.Close() }
