package screens

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	type point struct{ x, y float64 }
	var history [61]point
	position := point{384, 268}
	for i := range history {
		history[i] = position
	}
	palette := aHiddenPalette()
	whichColor, head := 0, 0
	s.input = func(in Input) {
		if in.PointerDown || in.PointerX != 0 || in.PointerY != 0 {
			position = point{in.PointerX, in.PointerY}
		}
	}
	s.render = func() {
		clearBlack(s.Canvas)
		history[head] = position
		vector.FillRect(s.Canvas, 388, 152, 310, 32, palette[whichColor], false)
		whichColor = (whichColor + 1) % len(palette)
		s.draw(s.Canvas, main, 64, 68)
		for i := 3; i >= 0; i-- {
			p := history[(head-i*20+len(history))%len(history)]
			s.draw(s.Canvas, pointers[i], p.x, p.y)
		}
		vector.FillRect(s.Canvas, float32(position.x), float32(position.y), 4, 2, palette[whichColor], false)
		vector.FillRect(s.Canvas, float32(position.x), float32(position.y), 2, 4, palette[whichColor], false)
		head = (head + 1) % len(history)
		vector.FillRect(s.Canvas, 0, 0, 768, 68, color.Black, false)
		vector.FillRect(s.Canvas, 0, 468, 768, 68, color.Black, false)
		vector.FillRect(s.Canvas, 0, 0, 64, 536, color.Black, false)
		vector.FillRect(s.Canvas, 704, 0, 64, 536, color.Black, false)
	}
}
