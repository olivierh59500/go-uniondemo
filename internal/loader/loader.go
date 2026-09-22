// Package loader presents the animated credits between the hall and a screen.
package loader

import (
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/assets"
	"github.com/olivierh59500/democonstructionkit/composite"
)

// Duration follows the embedded transition recording.
const Duration = 3777120 * time.Microsecond

type Screen struct {
	store       *assets.Store
	font, paper *ebiten.Image
	lines       []string
	tick        int
	rate        int
}

func New(id string, files fs.FS, rate int) (*Screen, error) {
	lines, ok := credits[id]
	if !ok {
		return nil, fmt.Errorf("loader: unknown screen %q", id)
	}
	s := &Screen{store: assets.New(files), lines: lines, rate: rate}
	var err error
	s.font, err = s.store.Texture("loader/loader.png")
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
	// Letters rise column by column, starting at the last row. Frame timing
	// stays independent of display refresh and keyboard repeat settings.
	columns := 20
	if len(s.lines) > 0 {
		columns = len(s.lines[0])
	}
	clock := float64(s.tick) * 140 * 60 / float64(s.rate)
	index := 0
	for col := 0; col < columns; col++ {
		for row := 0; row < len(s.lines); row++ {
			line := s.lines[len(s.lines)-1-row]
			var ch byte = ' '
			if col < len(line) {
				ch = line[col]
			}
			progress := math.Max(0, math.Min(1, (clock-float64(index*30))/50))
			y := 500 + (float64(374-row*16)-500)*progress
			if ch >= 32 {
				glyph := int(ch) - 32
				atlasColumns := s.font.Bounds().Dx() / 16
				r := image.Rect((glyph%atlasColumns)*16, (glyph/atlasColumns)*16, (glyph%atlasColumns+1)*16, (glyph/atlasColumns+1)*16)
				op := ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(154+col*16), y-8)
				composite.Instance{Image: s.font, Source: &r, Options: op}.Draw(s.paper)
			}
			index++
		}
	}
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(64, 60)
	dst.DrawImage(s.paper, &op)
}
func (s *Screen) Close() error { s.paper.Deallocate(); return s.store.Close() }
