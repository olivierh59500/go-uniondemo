package screens

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
)

func init() {
	factories["beatdis"] = (*Scene).beatDis
	factories["wow"] = (*Scene).wowScroller
}

func (s *Scene) beatDis() {
	font, backdrop := s.image("font.png"), s.image("backdrop.png")
	paper, scrollBack := s.image("wallpaper.png"), s.image("scroll.png")
	var letters [8]*ebiten.Image
	for i, name := range []string{"N", "O", "I", "N", "U", "E", "H", "T"} {
		letters[i] = s.image(name + ".png")
	}
	if s.err != nil {
		return
	}
	stage, scroll := s.surface(640, 297), s.surface(576, 100)
	const text = "     THE CAREBEARS PRESENT BEAT DIS!   A SCREEN FROM THE UNION DEMO.   THE UNION BRINGS TOGETHER THE CAREBEARS, THE EXCEPTIONS, THE REPLICANTS, THE TNT CREW, DELTA FORCE AND LEVEL 16.   ENJOY THE MUSIC AND THE MOVING LETTERS...                         "
	ring := s.ring(scroll, font, "union-beatdis", text, 5)
	if s.err != nil {
		return
	}
	wallpaper, err := composite.NewBackground(composite.BackgroundConfig{Source: image.Rect(0, 0, 640, 800), PeriodY: 800, Filter: ebiten.FilterNearest})
	if err != nil {
		s.err = err
		return
	}
	pattern, err := composite.NewBackground(composite.BackgroundConfig{Source: image.Rect(0, 0, 560, 32), PeriodX: 560, Filter: ebiten.FilterNearest})
	if err != nil {
		s.err = err
		return
	}
	phase := [8]float64{.2, .4, .6, .8, 1, 1.4, 1.6, 1.8}
	paperY, scrollX, xPhase := 328.0, 0.0, 0.0
	s.render = func() {
		clearBlack(s.Canvas)
		scroll.Clear()
		stage.Fill(color.RGBA{R: 160, A: 255})
		s.draw(s.Canvas, backdrop, 0, 0)
		if paperY > 0 {
			// Preserve the wallpaper's initial entrance before its first full repeat.
			s.draw(stage, paper, 0, paperY)
		} else {
			wallpaper.DrawAt(stage, paper, 0, paperY)
		}
		paperY -= 3
		if paperY <= -800 {
			paperY = 0
		}
		s.draw(s.Canvas, stage, 64, 60)
		xOffset := 360 + 180*math.Cos(xPhase/60)
		xPhase += .8
		for i, letter := range letters {
			phase[i] += .03
			s.draw(s.Canvas, letter, xOffset+100*math.Sin(phase[i]), 186+84*math.Cos(phase[i]*1.5))
		}
		pattern.DrawAt(scroll, scrollBack, scrollX, 34)
		scrollX -= 3
		if scrollX <= -560 {
			scrollX = 0
		}
		ring.Step()
		ring.DrawAt(scroll, 0, 0)
		s.draw(s.Canvas, scroll, 96, 357)
	}
}

func (s *Scene) wowScroller() {
	front, back, raster, font := s.image("mainF.png"), s.image("mainB.png"), s.image("raster.png"), s.image("fonts2.png")
	if s.err != nil {
		return
	}
	stage, scroll := s.surface(640, 400), s.surface(320, 200)
	const text = "   THE CAREBEARS PRESENT THE WOW SCROLLER!   A GIANT SCROLL THROUGH THE COLOURFUL WORLD OF THE UNION DEMO...   GREETINGS TO ALL MEMBERS OF THE UNION AND EVERYONE KEEPING THE ATARI ST ALIVE.                          "
	ring := s.ring(scroll, font, "union-wow", text, 5)
	if s.err != nil {
		return
	}
	rasterFill, err := composite.NewRasterOverlay(composite.RasterOverlayConfig{
		Image: raster, ScaleX: 640, ScaleY: 1, Alpha: 1,
		VelocityY: -2, WrapY: &composite.RasterWrap{Boundary: -108, Restart: -24, Inclusive: true},
		Filter: ebiten.FilterNearest, Blend: ebiten.BlendSourceAtop,
	})
	if err != nil {
		s.err = err
		return
	}
	backY, frontY := 0.0, 0.0
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		scroll.Clear()
		s.draw(stage, back, 0, backY)
		backY -= 2
		if backY <= -417 {
			backY = -17
		}
		ring.Step()
		ring.DrawAt(scroll, 0, 10)
		rasterFill.Draw(scroll)
		rasterFill.Step()
		s.transform(stage, scroll, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
		s.draw(stage, front, 0, frontY)
		frontY -= 2
		if frontY <= -417 {
			frontY = -17
		}
		s.draw(s.Canvas, stage, 64, 60)
	}
}
