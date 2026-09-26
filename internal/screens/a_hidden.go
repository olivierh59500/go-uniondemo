package screens

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/sprites"
	"github.com/olivierh59500/democonstructionkit/timeline"
)

func init() { factories["hidden"] = (*Scene).hiddenScreen }

func aHiddenPalette() [67]color.RGBA {
	values := [...]string{
		"000000", "00CC00", "11BB00", "22AA00", "339900", "448800", "557700", "666600", "775500", "884400", "993300", "AA2200", "BB1100", "CC0000", "DD0000", "EE0000", "FF0000",
		"EE0011", "DD0022", "CC0033", "BB0044", "AA0055", "990066", "880077", "770088", "660099", "5500AA", "4400BB", "3300CC", "2200DD", "1100EE", "0000FF",
		"1100FF", "2200FF", "3300FF", "4400FF", "5500FF", "6600FF", "7700FF", "8800FF", "9900FF", "AA00FF", "BB00FF", "CC00FF", "DD00FF", "EE00FF", "FF00FF",
		"FFBB44", "FFCC33", "FFDD22", "FFEE11", "FFFF00", "EEFF11", "DDFF22", "CCEF33", "BBEF44", "AAEF55", "99DF66", "88DF77", "77DF88", "66CF99", "55CFAA", "44CFBB", "33BFCC", "22BFDD", "11BFEE", "00BFFF",
	}
	var palette [67]color.RGBA
	for i, value := range values {
		packed, _ := strconv.ParseUint(value, 16, 32)
		palette[i] = color.RGBA{R: uint8(packed >> 16), G: uint8(packed >> 8), B: uint8(packed), A: 255}
	}
	return palette
}

func (s *Scene) hiddenScreen() {
	main := s.image("main.png")
	pointers := [4]*ebiten.Image{s.image("pointer1.png"), s.image("pointer2.png"), s.image("pointer3.png"), s.image("pointer4.png")}
	if s.err != nil {
		return
	}
	trailConfig := presets.UnionHiddenTrail(pointers[:])
	position := trailConfig.Initial
	trail, err := sprites.NewDelayedTrail(trailConfig)
	if err != nil {
		s.err = err
		return
	}
	palette := aHiddenPalette()
	paletteClock, err := timeline.NewPacedIndex(presets.UnionHiddenPalette(len(palette)))
	if err != nil {
		s.err = err
		return
	}
	s.input = func(in Input) {
		if in.PointerDown || in.PointerX != 0 || in.PointerY != 0 {
			position = geometry.Vec2{X: in.PointerX, Y: in.PointerY}
		}
	}
	s.render = func() {
		clearBlack(s.Canvas)
		if err := trail.SetPosition(position); err != nil {
			s.err = err
			return
		}
		if err := trail.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		vector.FillRect(s.Canvas, 388, 152, 310, 32, palette[paletteClock.Current()], false)
		paletteClock.Step()
		s.draw(s.Canvas, main, 64, 68)
		trail.Draw(s.Canvas)
		vector.FillRect(s.Canvas, float32(position.X), float32(position.Y), 4, 2, palette[paletteClock.Current()], false)
		vector.FillRect(s.Canvas, float32(position.X), float32(position.Y), 2, 4, palette[paletteClock.Current()], false)
		vector.FillRect(s.Canvas, 0, 0, 768, 68, color.Black, false)
		vector.FillRect(s.Canvas, 0, 468, 768, 68, color.Black, false)
		vector.FillRect(s.Canvas, 0, 0, 64, 536, color.Black, false)
		vector.FillRect(s.Canvas, 704, 0, 64, 536, color.Black, false)
	}
}
