package screens

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func init() { factories["replicants"] = (*Scene).replicants }

func (s *Scene) replicants() {
	main, mask, rasters := s.image("overlay.png"), s.image("rastersOverlay.png"), s.image("rasters.png")
	blue, red := s.image("fontBlue.png"), s.image("fontRed.png")
	pink, green, brown := s.image("rastersPink.png"), s.image("rastersGreen.png"), s.image("rastersBrown.png")
	var glyphs [14]*ebiten.Image
	for i, name := range []string{"T", "H", "E", "Space", "R", "E", "P", "L", "I", "C", "A", "N", "T", "S"} {
		glyphs[i] = s.image("sprite" + name + ".png")
	}
	if s.err != nil {
		return
	}
	stage, rasterStage := s.surface(640, 400), s.surface(640, 200)
	rasterFill, err := composite.NewRasterOverlay(composite.RasterOverlayConfig{
		Image: rasters, Y: -53, ScaleX: 1, ScaleY: 1, Alpha: 1,
		VelocityY: 1, WrapY: &composite.RasterWrap{Boundary: -53, Restart: -221, Inclusive: true},
		Filter: ebiten.FilterNearest, Blend: ebiten.BlendSourceAtop,
	})
	if err != nil {
		s.err = err
		return
	}
	const text = "                    THE REPLICANTS PRESENT THEIR WOBBLY SPRITE SCREEN!   WELCOME TO THE UNION DEMO.   GREETINGS TO THE CAREBEARS, THE EXCEPTIONS, THE TNT CREW, DELTA FORCE AND LEVEL 16.   ENJOY THE RASTERS, THE MUSIC AND THE DANCING LETTERS!                           "
	top, err := s.bitmap(red, "union-replicants").Scrolling(text)
	if err != nil {
		s.err = err
		return
	}
	bottom, err := s.bitmap(blue, "union-replicants").Scrolling(text)
	if err != nil {
		s.err = err
		return
	}
	if s.err != nil {
		return
	}
	topY, bottomY := [3]float64{94, 124, 154}, [3]float64{204, 234, 264}
	topDY, bottomDY := [3]float64{2, 2, 2}, [3]float64{-2, -2, -2}
	topRasters, bottomRasters := [3]*ebiten.Image{pink, green, brown}, [3]*ebiten.Image{brown, green, pink}
	formation, err := replicantsFormation()
	if err != nil {
		s.err = err
		return
	}
	letters, err := sprites.NewGroup(sprites.GroupConfig{
		Frames: glyphs[:], Count: len(glyphs), FrameStride: 1,
		Formation: formation.At, Speed: 1, Filter: ebiten.FilterNearest,
	})
	if err != nil {
		s.err = err
		return
	}
	scrollX, scrollSpeed := -640.0, 6.0
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
		rasterFill.Step()
		s.draw(rasterStage, mask, 0, 0)
		rasterFill.Draw(rasterStage)
		s.draw(s.Canvas, rasterStage, 64, 206)
		if err := letters.Update(kit.Frame{Time: float64(s.Frame) / 60}); err != nil {
			s.err = err
			return
		}
		letters.Draw(s.Canvas)
	}
}
