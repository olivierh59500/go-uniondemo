package screens

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func init() {
	factories["tnt2"] = buildTNT2
	factories["starballs"] = buildStarballs
	factories["level16"] = buildLevel16
}

func buildTNT2(s *Scene) {
	stage := s.surface(640, 400)
	layers := []*ebiten.Image{s.image("blueLayer.png"), s.image("brownLayer.png"), s.image("greenLayer.png")}
	front, back := s.image("overlay.png"), s.image("overlay2.png")
	font := s.bitmap(s.image("fonts.png"), "union-tnt2")
	text := "                     THE TNT CREW PRESENTS THE SUPERSCROLLER! THREE INDEPENDENT BACKGROUND LAYERS AND A GIANT SCROLLINE FOR THE UNION DEMO. GREETINGS TO TEX, THE CAREBEARS, THE REPLICANTS, DELTA FORCE AND LEVEL 16. USE THE ARROW KEYS TO CHANGE THE SCROLL SPEED AND DIRECTION. MUSIC BY MAD MAX.                           "
	textView, err := scrolling.NewBitmapText(font, text, 0)
	if err != nil {
		s.err = err
		return
	}
	backgrounds, err := composite.NewBackground(composite.BackgroundConfig{PeriodX: 640, Filter: ebiten.FilterNearest})
	if err != nil {
		s.err = err
		return
	}
	parallax, err := motion.NewWrapBank(motion.WrapBankConfig{
		Start: []float64{-640, -640, -640}, Velocity: []float64{-2, -4, -6},
		Lower: &motion.WrapLimit{Boundary: -640, Restart: 0},
		Upper: &motion.WrapLimit{Boundary: 0, Restart: -640},
	})
	if err != nil {
		s.err = err
		return
	}
	scrollClock, err := motion.NewWrapBank(presets.UnionTNT2TextWrap(textView.Width()))
	if err != nil {
		s.err = err
		return
	}
	scrollSpeed := 4.0
	direction, held := -1.0, false
	s.input = func(in Input) {
		if in.Left {
			direction = -1
		}
		if in.Right {
			direction = 1
		}
		switch in.Number {
		case 1:
			if err := parallax.SetVelocity(0, direction*2); err != nil {
				s.err = err
			}
		case 2:
			if err := parallax.AddVelocity(1, direction); err != nil {
				s.err = err
			}
		case 3:
			if err := parallax.AddVelocity(2, direction); err != nil {
				s.err = err
			}
		case 4:
			scrollSpeed = max(0, min(9, scrollSpeed-direction))
		}
		if in.Action && !held {
			parallax.ReverseAll()
		}
		held = in.Action
		if in.Up {
			scrollSpeed = math.Min(9, scrollSpeed+.05)
		}
		if in.Down {
			scrollSpeed = math.Max(0, scrollSpeed-.05)
		}
		if err := scrollClock.SetVelocity(0, -scrollSpeed); err != nil {
			s.err = err
		}
	}
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		s.draw(stage, back, 0, 0)
		parallax.Step()
		for i, layer := range layers {
			backgrounds.DrawAt(stage, layer, parallax.At(i), 0)
		}
		if err := textView.DrawWindow(stage, scrollClock.At(0), 180, 640); err != nil {
			s.err = err
			return
		}
		scrollClock.Step()
		s.draw(stage, front, 0, 0)
		s.draw(s.Canvas, stage, 64, 60)
	}
}

func buildStarballs(s *Scene) {
	font, bob1, bob2, logo := s.image("font.png"), s.image("union_bob1.png"), s.image("union_bob2.png"), s.image("union_logo.png")
	text := " THE TNT CREW PRESENTS STARBALLS, A SCREEN FROM THE UNION DEMO! WATCH THE BALLS CHANGE COLOUR AS THEY CROSS THE LOGO AND THE SCROLLINE. USE UP AND DOWN TO CHANGE THE NUMBER OF STARBALLS. GREETINGS TO ALL MEMBERS OF THE UNION! MUSIC BY MAD MAX.   "
	scroll, err := scrolling.NewRing(scrolling.RingConfig{
		Text: text, Font: s.bitmap(font, "union-starballs"), Viewport: 320, Speed: 3, Controls: true,
	})
	if err != nil {
		s.err = err
		return
	}
	config, err := presets.UnionStarballs(bob1, bob2, s.rnd, func(mask *ebiten.Image) {
		s.draw(mask, logo, -32, -34)
		scroll.DrawAt(mask, 0, 188)
	})
	if err != nil {
		s.err = err
		return
	}
	layer, err := sprites.NewMaskedProjectedField(config)
	if err != nil {
		s.err = err
		return
	}
	s.closers = append(s.closers, layer.Close)
	counts := []int{32, 40, 48, 64, 72, 80, 88, 96, 112, 128}
	selected, held := 0, false
	s.input = func(in Input) {
		pressed := in.Up || in.Down || in.Action
		if in.Number >= 1 && in.Number <= 10 || pressed && !held {
			if in.Number >= 1 && in.Number <= 10 {
				selected = in.Number - 1
			} else if in.Down {
				selected = (selected + len(counts) - 1) % len(counts)
			} else {
				selected = (selected + 1) % len(counts)
			}
			s.err = layer.Field().ResetCount(counts[selected])
		}
		held = pressed
	}
	s.render = func() {
		s.Canvas.Fill(color.RGBA{B: 64, A: 255})
		scroll.Step()
		if err := layer.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		layer.Draw(s.Canvas)
	}
}

func buildLevel16(s *Scene) {
	back, bob := s.image("l16bck.png"), s.image("l16bob.png")
	font, water, raster := s.image("l16font.png"), s.image("l16water.png"), s.image("l16raster.png")
	grid := s.bitmap(font, "union-level16")
	// The same recycled glyph positions drive a vertical scroll without rotating its letters.
	text := "LEVEL 16 PRESENTS THE FULLSCREEN! ANOTHER SCREEN FROM THE GREAT UNION DEMO. GREETINGS TO ALL OUR FRIENDS IN THE UNION: TEX, THE CAREBEARS, TNT CREW, DELTA FORCE AND THE REPLICANTS. ENJOY THE WATER, RASTERS AND THE BOUNCING BALL!   "
	vertical, err := scrolling.New(scrolling.Config{X: 698, Recycled: &scrolling.RecycledConfig{Vertical: true, Ring: scrolling.RingConfig{Text: text, Font: grid, Viewport: 536, Speed: 2}}})
	if err != nil {
		s.err = err
		return
	}
	s.closers = append(s.closers, vertical.Close)
	waterMotion, err := motion.NewWrapBank(motion.WrapBankConfig{
		Start: []float64{0}, Velocity: []float64{2},
		Upper: &motion.WrapLimit{Boundary: 220, Restart: 0, Inclusive: true},
	})
	if err != nil {
		s.err = err
		return
	}
	rasterMotion, err := motion.NewWrapBank(motion.WrapBankConfig{
		Start: []float64{120}, Velocity: []float64{-2},
		Lower: &motion.WrapLimit{Boundary: -25, Restart: 120, Inclusive: true},
	})
	if err != nil {
		s.err = err
		return
	}
	phase := 0.0
	orbit := motion.DefaultNestedOrbit(motion.Point{X: 384, Y: 268}, motion.Point{X: 192, Y: 536 / 2.7})
	s.render = func() {
		clearBlack(s.Canvas)
		s.err = vertical.Update(kit.Frame{Tick: s.Frame})
		vertical.Draw(s.Canvas)
		s.draw(s.Canvas, water, 20, waterMotion.At(0))
		waterMotion.Step()
		s.draw(s.Canvas, raster, 300, rasterMotion.At(0))
		rasterMotion.Step()
		s.draw(s.Canvas, back, 0, 0)
		phase += .008
		position := orbit.At(phase)
		x, y := position.X, position.Y
		s.transform(s.Canvas, bob, x, y, 1, 1, 0, float64(bob.Bounds().Dx())/2, float64(bob.Bounds().Dy())/2, 1, ebiten.BlendSourceOver)
	}
}

// Keep rectangular sampling explicit for the atlas-based effects in this group.
func unionRegion(x, y, w, h float64) composite.Region {
	return composite.Region{X: x, Y: y, Width: w, Height: h}
}
