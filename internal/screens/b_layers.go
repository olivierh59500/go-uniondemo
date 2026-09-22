package screens

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

func init() {
	factories["tnt2"] = buildTNT2
	factories["starballs"] = buildStarballs
	factories["level16"] = buildLevel16
}

func unionBitmap(atlas *ebiten.Image, w, h float64) scrolling.BitmapGrid {
	return scrolling.BitmapGrid{Image: atlas, Width: w, Height: h, Columns: max(1, int(float64(atlas.Bounds().Dx())/w)), ColumnSpan: float64(atlas.Bounds().Dx()) / w, First: 32, Filter: ebiten.FilterNearest}
}

func buildTNT2(s *Scene) {
	stage := s.surface(640, 400)
	layers := []*ebiten.Image{s.image("blueLayer.png"), s.image("brownLayer.png"), s.image("greenLayer.png")}
	front, back := s.image("overlay.png"), s.image("overlay2.png")
	font := unionBitmap(s.image("fonts.png"), 64, 40)
	text := "                     THE TNT CREW PRESENTS THE SUPERSCROLLER! THREE INDEPENDENT BACKGROUND LAYERS AND A GIANT SCROLLINE FOR THE UNION DEMO. GREETINGS TO TEX, THE CAREBEARS, THE REPLICANTS, DELTA FORCE AND LEVEL 16. USE THE ARROW KEYS TO CHANGE THE SCROLL SPEED AND DIRECTION. MUSIC BY MAD MAX.                           "
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
			s.draw(stage, layer, positions[i], 0)
			s.draw(stage, layer, positions[i]+640, 0)
		}
		// Print only intersecting glyphs; the message never needs a giant texture.
		first := max(0, int(-scrollX/64))
		last := min(len(text), first+12)
		font.Print(stage, text[first:last], scrollX+float64(first)*64, 180, 1, 1)
		scrollX -= scrollSpeed
		if scrollX < -float64(len(text)*64-640) {
			scrollX = -640
		}
		s.draw(stage, front, 0, 0)
		s.draw(s.Canvas, stage, 64, 60)
	}
}

func buildStarballs(s *Scene) {
	type star struct{ x, y, z float64 }
	base, mask := s.surface(320, 200), s.surface(320, 200)
	font, bob1, bob2, logo := s.image("font.png"), s.image("union_bob1.png"), s.image("union_bob2.png"), s.image("union_logo.png")
	text := " THE TNT CREW PRESENTS STARBALLS, A SCREEN FROM THE UNION DEMO! WATCH THE BALLS CHANGE COLOUR AS THEY CROSS THE LOGO AND THE SCROLLINE. USE UP AND DOWN TO CHANGE THE NUMBER OF STARBALLS. GREETINGS TO ALL MEMBERS OF THE UNION! MUSIC BY MAD MAX.   "
	scroll := s.ring(mask, font, 15, 8, 32, text, 3)
	stars := make([]star, 32)
	randomize := func(p *star) {
		p.x = math.Floor(s.rnd()*49) - 25
		p.y = math.Floor(s.rnd()*49) - 25
		p.z = math.Floor(s.rnd()*30) + 1
	}
	for i := range stars {
		randomize(&stars[i])
	}
	counts := []int{32, 40, 48, 64, 72, 80, 88, 96, 112, 128}
	selected, held := 0, false
	camera := geometry.Camera{Center: geometry.Vec2{X: 160, Y: 100}, Focal: 64, Near: .001}
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
			stars = make([]star, counts[selected])
			for i := range stars {
				randomize(&stars[i])
			}
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
		for i := range stars {
			p := &stars[i]
			p.z -= .2
			if p.z <= 0 {
				randomize(p)
				p.z = 32
			}
			projected, _, visible := camera.Project(geometry.Vec3{X: p.x, Y: p.y, Z: p.z})
			if !visible {
				continue
			}
			x, y := projected.X, projected.Y
			if x < 0 || x > 320 || y < 0 || y > 200 {
				continue
			}
			size := (1 - p.z/32) * 5 / 8
			alpha := math.Floor((1-p.z/32)*255) / 255
			s.transform(base, bob1, x, y, size, size, 0, 0, 0, alpha, ebiten.BlendSourceOver)
			s.transform(mask, bob2, x, y, size, size, 0, 0, 0, alpha, ebiten.BlendSourceAtop)
		}
		s.transform(s.Canvas, base, 64, 68, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
		s.transform(s.Canvas, mask, 64, 68, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
	}
}

func buildLevel16(s *Scene) {
	back, bob := s.image("l16bck.png"), s.image("l16bob.png")
	font, water, raster := s.image("l16font.png"), s.image("l16water.png"), s.image("l16raster.png")
	grid := unionBitmap(font, 32, 32)
	// The same recycled glyph positions drive a vertical scroll without rotating its letters.
	text := "LEVEL 16 PRESENTS THE FULLSCREEN! ANOTHER SCREEN FROM THE GREAT UNION DEMO. GREETINGS TO ALL OUR FRIENDS IN THE UNION: TEX, THE CAREBEARS, TNT CREW, DELTA FORCE AND THE REPLICANTS. ENJOY THE WATER, RASTERS AND THE BOUNCING BALL!   "
	vertical, err := scrolling.NewRing(scrolling.RingConfig{Text: text, Font: grid, Viewport: 536, Speed: 2})
	if err != nil {
		s.err = err
		return
	}
	waterY, rasterY, phase := 0.0, 120.0, 0.0
	s.render = func() {
		clearBlack(s.Canvas)
		vertical.Step()
		for _, letter := range vertical.Letters() {
			if region, ok := grid.Region(letter.Rune); ok {
				s.part(s.Canvas, font, region, 698, letter.X, 1, 1)
			}
		}
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
		x := 384 + 192*math.Cos(phase*4-math.Cos(phase-.1))
		y := 268 + 536/2.7*-math.Sin(phase*2.3-math.Cos(phase-.1))
		s.transform(s.Canvas, bob, x, y, 1, 1, 0, float64(bob.Bounds().Dx())/2, float64(bob.Bounds().Dy())/2, 1, ebiten.BlendSourceOver)
	}
}

// Keep rectangular sampling explicit for the atlas-based effects in this group.
func unionRegion(x, y, w, h float64) composite.Region {
	return composite.Region{X: x, Y: y, Width: w, Height: h}
}
