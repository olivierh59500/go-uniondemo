package screens

import (
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/assets"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

// Input is shared by keyboard, mouse and touch hosts.
type Input struct {
	Left, Right, Up, Down, Action bool
	Number                        int // 1..9, with 10 representing the zero key.
	PointerX, PointerY            float64
	PointerDown                   bool
}

type Scene struct {
	Descriptor   Descriptor
	Canvas       *ebiten.Image
	Frame        uint64
	Controls     Input
	VoiceVolumes [3]uint8
	Music        string
	store        *assets.Store
	surfaces     []*ebiten.Image
	render       func()
	input        func(Input)
	closers      []func() error
	err          error
	random       uint32
}

var factories = map[string]func(*Scene){}

func New(id string, files fs.FS) (*Scene, error) {
	d, ok := Find(id)
	if !ok {
		return nil, fmt.Errorf("unknown Union screen %q", id)
	}
	build, ok := factories[id]
	if !ok {
		return nil, fmt.Errorf("Union screen %q is not available", id)
	}
	s := &Scene{Descriptor: d, Music: d.Music, store: assets.New(files), random: 42}
	s.Canvas = s.surface(d.Width, d.Height)
	build(s)
	if s.err != nil {
		s.Close()
		return nil, s.err
	}
	if s.render == nil {
		s.Close()
		return nil, fmt.Errorf("screen %q has no renderer", id)
	}
	s.render()
	return s, nil
}

func Available(id string) bool { return factories[id] != nil }
func (s *Scene) Update(in Input) error {
	s.Controls = in
	if s.input != nil {
		s.input(in)
	}
	s.Frame++
	s.render()
	return s.err
}
func (s *Scene) Draw(dst *ebiten.Image) { dst.DrawImage(s.Canvas, nil) }
func (s *Scene) Close() error {
	for _, close := range s.closers {
		close()
	}
	s.closers = nil
	for _, img := range s.surfaces {
		img.Deallocate()
	}
	s.surfaces = nil
	return s.store.Close()
}
func (s *Scene) image(name string) *ebiten.Image {
	img, err := s.store.Texture(s.Descriptor.Directory + "/" + name)
	if err != nil {
		s.err = err
		return s.surface(1, 1)
	}
	return img
}
func (s *Scene) surface(w, h int) *ebiten.Image {
	img := ebiten.NewImageWithOptions(image.Rect(0, 0, w, h), &ebiten.NewImageOptions{Unmanaged: true})
	s.surfaces = append(s.surfaces, img)
	return img
}
func (s *Scene) bitmap(atlas *ebiten.Image, recipe string) scrolling.BitmapGrid {
	grid, err := presets.BitmapFont(recipe, atlas, ebiten.FilterNearest)
	if err != nil {
		s.err = err
	}
	return grid
}
func (s *Scene) ring(dst, atlas *ebiten.Image, recipe, text string, speed float64) *scrolling.Ring {
	r, err := scrolling.NewRing(scrolling.RingConfig{Text: text, Font: s.bitmap(atlas, recipe), Viewport: float64(dst.Bounds().Dx()), Speed: speed, Controls: true})
	if err != nil {
		s.err = err
	}
	return r
}

func (s *Scene) draw(dst, src *ebiten.Image, x, y float64) {
	s.transform(dst, src, x, y, 1, 1, 0, 0, 0, 1, ebiten.BlendSourceOver)
}
func (s *Scene) transform(dst, src *ebiten.Image, x, y, sx, sy, angle, hx, hy, alpha float64, blend ebiten.Blend) {
	op := ebiten.DrawImageOptions{Filter: ebiten.FilterNearest, Blend: blend}
	op.GeoM.Translate(-hx, -hy)
	op.GeoM.Scale(sx, sy)
	op.GeoM.Rotate(angle * math.Pi / 180)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleAlpha(float32(alpha))
	composite.Instance{Image: src, Options: op}.Draw(dst)
}
func (s *Scene) part(dst, src *ebiten.Image, r composite.Region, x, y, sx, sy float64) {
	op := ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
	op.GeoM.Scale(sx, sy)
	op.GeoM.Translate(x, y)
	composite.DrawRegion(dst, src, r, &op)
}
func (s *Scene) rnd() float64 {
	s.random = s.random*1664525 + 1013904223
	return float64(s.random) / 4294967296
}
func clearBlack(dst *ebiten.Image) { dst.Fill(color.Black) }
