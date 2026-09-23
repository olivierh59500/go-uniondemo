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
)

// Duration follows the embedded transition recording.
const Duration = 3777120 * time.Microsecond

type Screen struct {
	store       *assets.Store
	font, paper *ebiten.Image
	reveal      *scrolling.Reveal
	tick        int
	rate        int
}

func New(id string, files fs.FS, rate int) (*Screen, error) {
	lines, ok := credits[id]
	if !ok {
		return nil, fmt.Errorf("loader: unknown screen %q", id)
	}
	s := &Screen{store: assets.New(files), rate: rate}
	var err error
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
func (s *Screen) Update() { s.tick++ }
func (s *Screen) Done() bool {
	return time.Duration(s.tick)*time.Second >= Duration*time.Duration(s.rate)
}
func (s *Screen) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	s.paper.Clear()
	// The shared reveal uses a fixed clock independent from display refresh.
	s.reveal.DrawAt(s.paper, float64(s.tick)*140*60/float64(s.rate))

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(64, 60)
	dst.DrawImage(s.paper, &op)
}
func (s *Screen) Close() error { s.paper.Deallocate(); return s.store.Close() }
