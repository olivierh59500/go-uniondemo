package screens

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/modulation"
)

func init() { factories["delta"] = (*Scene).deltaForce }

func (s *Scene) deltaForce() {
	slab, logos, gold, balls, font := s.image("slab.png"), s.image("logos.png"), s.image("gold_backdrop.png"), s.image("ball.png"), s.image("questfont.png")
	if s.err != nil {
		return
	}
	logoStage, merge, goldStage := s.surface(640, 130), s.surface(640, 480), s.surface(640, 472)
	scroll, word, wordMerge := s.surface(640, 150), s.surface(640, 150), s.surface(640, 480)
	grid := s.bitmap(font, "union-delta")
	grid.Print(word, "MEGA DEMO", 0, 0, 1, 1)
	const text = "DELTA FORCE PRESENTS THE ^P3SPHER-I-COOL SCREEN FROM THE UNION DEMO ^P3!!   THREE BOUNCING BALLS FOLLOW THE MUSIC WHILE THE GOLDEN SCROLL WAVES ACROSS THE SCREEN.   GREETINGS TO EVERY MEMBER OF THE UNION...                         "
	ring := s.ring(scroll, font, "union-delta", text, 6)
	if s.err != nil {
		return
	}
	// The introductory title and scrolling text share one continuous wave phase.
	wave := composite.WaveStrips{Axis: composite.Columns, Thickness: 1, Filter: ebiten.FilterLinear, Waves: []composite.StripWave{{Phase: 10, Amplitude: 50, Spatial: .002, Speed: -.03}}}
	logoScale, logoStep, goldY := 1.0, .02, 0.0
	tile, wordMode, wordWait := 0, 0, 200
	wordX, wordDX := -640.0, 6.0
	var voiceChanges [3]modulation.Change[uint8]
	var ballEnvelopes [3]*modulation.Decay
	for i := range ballEnvelopes {
		ballEnvelopes[i], _ = modulation.NewDecay(modulation.DecayConfig{Peak: 7, Rate: 1})
	}
	drawWord := func() {
		wave.DrawAt(wordMerge, word, 0, 320)
		wave.Advance()
		s.transform(wordMerge, goldStage, 0, 0, 1, 1, 0, 0, 0, 1, ebiten.BlendSourceAtop)
		s.draw(s.Canvas, wordMerge, 64+wordX, 28)
		vector.FillRect(s.Canvas, 0, 0, 64, 536, color.Black, false)
	}
	s.render = func() {
		clearBlack(s.Canvas)
		scroll.Clear()
		merge.Clear()
		logoStage.Clear()
		wordMerge.Clear()
		s.draw(s.Canvas, slab, 84, 218)
		columns := max(1, logos.Bounds().Dx()/640)
		s.part(logoStage, logos, composite.Region{X: float64(tile % columns * 640), Y: float64(tile / columns * 130), Width: 640, Height: 130}, 0, 0, 1, 1)
		s.transform(s.Canvas, logoStage, 384, 64, 1, logoScale, 0, 320, 65, 1, ebiten.BlendSourceOver)
		logoScale -= logoStep
		if logoScale <= 0 {
			logoStep = -.02
			tile = (tile + 1) % 2
		}
		if logoScale >= 1 {
			logoStep = .02
		}
		for i, volume := range s.VoiceVolumes {
			frame := int(ballEnvelopes[i].Step(voiceChanges[i].Sample(volume), 1))
			x := [3]float64{244, 355, 464}[i]
			s.part(s.Canvas, balls, composite.Region{X: float64(frame * 96), Width: 96, Height: 114}, x, 177, 1, 1)
		}
		s.draw(goldStage, gold, 0, 236-goldY)
		s.draw(goldStage, gold, 0, 354-goldY)
		goldY -= 3
		if goldY <= -59 {
			goldY = 0
		}
		if wordMode == 0 {
			wordX += wordDX
			if wordX >= 32 {
				wordDX, wordMode = -6, 1
			}
			drawWord()
		}
		if wordMode == 1 {
			wordWait--
			if wordWait <= 0 {
				wordMode = 2
			}
			drawWord()
		}
		if wordMode == 2 {
			wordX += wordDX
			if wordX <= -640 {
				wordDX, wordMode = 0, 3
			}
			drawWord()
		}
		if wordMode == 3 {
			ring.Step()
			ring.DrawAt(scroll, 0, 0)
			wave.DrawAt(merge, scroll, 0, 320)
			wave.Advance()
			s.transform(merge, goldStage, 0, 0, 1, 1, 0, 0, 0, 1, ebiten.BlendSourceAtop)
			s.draw(s.Canvas, merge, 64, 28)
		}
	}
}
