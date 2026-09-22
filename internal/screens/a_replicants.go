package screens

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

func init() { factories["replicants"] = (*Scene).replicants }

// aBitmapText retains a fixed-width text strip without allocating a giant image.
func (s *Scene) aBitmapText(atlas *ebiten.Image, text string, width, height int, first rune) *scrolling.Scrolling {
	columns := max(1, atlas.Bounds().Dx()/width)
	glyphs := make([]*ebiten.Image, len([]rune(text)))
	for i, ch := range []rune(text) {
		tile := int(ch - first)
		if tile < 0 {
			continue
		}
		x, y := tile%columns*width, tile/columns*height
		region := image.Rect(x, y, x+width, y+height)
		if region.In(atlas.Bounds()) {
			glyphs[i] = atlas.SubImage(region).(*ebiten.Image)
		}
	}
	textStrip, err := scrolling.FromImages(glyphs, float64(width))
	if err != nil {
		s.err = err
	}
	return textStrip
}

func (s *Scene) replicants() {
	main, mask, rasters := s.image("overlay.png"), s.image("rastersOverlay.png"), s.image("rasters.png")
	blue, red := s.image("fontBlue.png"), s.image("fontRed.png")
	pink, green, brown := s.image("rastersPink.png"), s.image("rastersGreen.png"), s.image("rastersBrown.png")
	var sprites [14]*ebiten.Image
	for i, name := range []string{"T", "H", "E", "Space", "R", "E", "P", "L", "I", "C", "A", "N", "T", "S"} {
		sprites[i] = s.image("sprite" + name + ".png")
	}
	if s.err != nil {
		return
	}
	stage, rasterStage := s.surface(640, 400), s.surface(640, 200)
	const text = "                    THE REPLICANTS PRESENT THEIR WOBBLY SPRITE SCREEN!   WELCOME TO THE UNION DEMO.   GREETINGS TO THE CAREBEARS, THE EXCEPTIONS, THE TNT CREW, DELTA FORCE AND LEVEL 16.   ENJOY THE RASTERS, THE MUSIC AND THE DANCING LETTERS!                           "
	top, bottom := s.aBitmapText(red, text, 64, 64, 32), s.aBitmapText(blue, text, 64, 64, 32)
	if s.err != nil {
		return
	}
	topY, bottomY := [3]float64{94, 124, 154}, [3]float64{204, 234, 264}
	topDY, bottomDY := [3]float64{2, 2, 2}, [3]float64{-2, -2, -2}
	topRasters, bottomRasters := [3]*ebiten.Image{pink, green, brown}, [3]*ebiten.Image{brown, green, pink}
	phase := [14]float64{}
	for i := range phase {
		phase[i] = float64(i) * .5
	}
	scrollX, scrollSpeed, rasterY := -640.0, 6.0, -53.0
	previous := Input{}
	s.input = func(in Input) {
		if in.Right && !previous.Right {
			scrollSpeed = math.Min(scrollSpeed+1, 16)
		}
		if in.Left && !previous.Left {
			scrollSpeed = math.Max(scrollSpeed-1, 0)
		}
		previous = in
	}
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		rasterStage.Clear()
		gold := color.RGBA{R: 134, G: 100, A: 255}
		vector.FillRect(stage, 0, 0, 640, 80, gold, false)
		vector.FillRect(stage, 0, 326, 640, 76, gold, false)
		state := scrolling.IdentityState()
		state.X, state.Y = scrollX, 14
		state.First = max(0, int(-scrollX/64)-1)
		state.End = min(len(text), state.First+13)
		top.DrawAt(stage, state)
		state.Y = 326
		bottom.DrawAt(stage, state)
		scrollX -= scrollSpeed
		if scrollX < -float64(len(text)*64-640) {
			scrollX = -640
		}
		s.draw(stage, main, 0, 0)
		for i := range topY {
			topY[i] += topDY[i]
			if topY[i] < 94 || topY[i] > 160 {
				topDY[i] = -topDY[i]
			}
			bottomY[i] += bottomDY[i]
			if bottomY[i] < 204 || bottomY[i] > 270 {
				bottomDY[i] = -bottomDY[i]
			}
		}
		for i, raster := range topRasters {
			s.transform(s.Canvas, raster, 0, 60+topY[i], 390, 1, 0, 0, 0, 1, ebiten.BlendSourceOver)
		}
		for i, raster := range bottomRasters {
			s.transform(s.Canvas, raster, 0, 60+bottomY[i], 390, 1, 0, 0, 0, 1, ebiten.BlendSourceOver)
		}
		s.draw(s.Canvas, stage, 64, 60)
		rasterY++
		if rasterY >= -53 {
			rasterY = -221
		}
		s.draw(rasterStage, mask, 0, 0)
		s.transform(rasterStage, rasters, 0, rasterY, 1, 1, 0, 0, 0, 1, ebiten.BlendSourceAtop)
		s.draw(s.Canvas, rasterStage, 64, 206)
		for i, sprite := range sprites {
			phase[i] += .08
			s.draw(s.Canvas, sprite, 64+40+float64(i)*40+30*math.Cos(phase[i]), 60+180+30*math.Sin(phase[i]))
		}
	}
}
