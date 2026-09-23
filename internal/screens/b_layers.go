package screens

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/motion"
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
	positions := [3]float64{-640, -640, -640}
	speeds := [3]float64{-2, -4, -6}
	scrollX, scrollSpeed := -640.0, 4.0
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
			speeds[0] = direction * 2
		case 2:
			speeds[1] += direction
		case 3:
			speeds[2] += direction
		case 4:
			scrollSpeed = max(0, min(9, scrollSpeed-direction))
		}
		if in.Action && !held {
			for i := range speeds {
				speeds[i] = -speeds[i]
			}
		}
		held = in.Action
		if in.Up {
			scrollSpeed = math.Min(9, scrollSpeed+.05)
		}
		if in.Down {
			scrollSpeed = math.Max(0, scrollSpeed-.05)
		}
	}
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		s.draw(stage, back, 0, 0)
		for i, layer := range layers {
			positions[i] += speeds[i]
			if positions[i] > 0 {
				positions[i] = -640
			}
			if positions[i] < -640 {
				positions[i] = 0
			}
			backgrounds.DrawAt(stage, layer, positions[i], 0)
		}
		if err := textView.DrawWindow(stage, scrollX, 180, 640); err != nil {
			s.err = err
			return
		}
		scrollX -= scrollSpeed
		if scrollX < -(textView.Width() - 640) {
			scrollX = -640
		}
		s.draw(stage, front, 0, 0)
		s.draw(s.Canvas, stage, 64, 60)
	}
}

func buildStarballs(s *Scene) {
	base, mask := s.surface(320, 200), s.surface(320, 200)
	font, bob1, bob2, logo := s.image("font.png"), s.image("union_bob1.png"), s.image("union_bob2.png"), s.image("union_logo.png")
	text := " THE TNT CREW PRESENTS STARBALLS, A SCREEN FROM THE UNION DEMO! WATCH THE BALLS CHANGE COLOUR AS THEY CROSS THE LOGO AND THE SCROLLINE. USE UP AND DOWN TO CHANGE THE NUMBER OF STARBALLS. GREETINGS TO ALL MEMBERS OF THE UNION! MUSIC BY MAD MAX.   "
	scroll := s.ring(mask, font, "union-starballs", text, 3)
	field, err := sprites.NewProjectedField(sprites.ProjectedFieldConfig{
		Field: sprites.FieldConfig{Count: 32, Depth: sprites.DepthRespawn, Near: 0, Far: 32,
			Spawn: func(_ int, reset bool) sprites.Point {
				p := sprites.Point{X: math.Floor(s.rnd()*49) - 25, Y: math.Floor(s.rnd()*49) - 25, Z: math.Floor(s.rnd()*30) + 1}
				if reset {
					p.Z = 32
				}
				return p
			}},
		View:     sprites.FieldView{Camera: geometry.Camera{Center: geometry.Vec2{X: 160, Y: 100}, Focal: 64, Near: .001}},
		Velocity: geometry.Vec3{Z: -.2}, Delta: 1, RendererCapacity: 128,
	})
	if err != nil {
		s.err = err
		return
	}
	s.closers = append(s.closers, field.Close)
	style := sprites.FieldStyle{DrawImages: true, Sample: func(p sprites.FieldSample, a *sprites.FieldAppearance) bool {
		if p.X < 0 || p.X > 320 || p.Y < 0 || p.Y > 200 {
			return false
		}
		size := (1 - p.Z/32) * 5 / 8
		a.ScaleX, a.ScaleY = size, size
		a.Tint.ScaleAlpha(float32(math.Floor((1-p.Z/32)*255) / 255))
		return true
	}}
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
			s.err = field.ResetCount(counts[selected])
		}
		held = pressed
	}
	s.render = func() {
		s.Canvas.Fill(color.RGBA{B: 64, A: 255})
		base.Clear()
		mask.Clear()
		s.draw(mask, logo, -32, -34)
		scroll.Step()
		scroll.DrawAt(mask, 0, 188)
		if err := field.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		style.Image, style.Blend = bob1, ebiten.BlendSourceOver
		field.DrawStyle(base, style)
		style.Image, style.Blend = bob2, ebiten.BlendSourceAtop
		field.DrawStyle(mask, style)
		s.transform(s.Canvas, base, 64, 68, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
		s.transform(s.Canvas, mask, 64, 68, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
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
	waterY, rasterY, phase := 0.0, 120.0, 0.0
	orbit := motion.DefaultNestedOrbit(motion.Point{X: 384, Y: 268}, motion.Point{X: 192, Y: 536 / 2.7})
	s.render = func() {
		clearBlack(s.Canvas)
		s.err = vertical.Update(kit.Frame{Tick: s.Frame})
		vertical.Draw(s.Canvas)
		s.draw(s.Canvas, water, 20, waterY)
		waterY += 2
		if waterY >= 220 {
			waterY = 0
		}
		s.draw(s.Canvas, raster, 300, rasterY)
		rasterY -= 2
		if rasterY <= -25 {
			rasterY = 120
		}
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
