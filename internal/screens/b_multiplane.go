package screens

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

func init() { factories["multiplane"] = buildMultiplane }

func buildMultiplane(s *Scene) {
	rasters, mountains, logo, font := s.image("rast.png"), s.image("mountains.png"), s.image("logo.png"), s.image("bgfont.png")
	background, stage := s.surface(640, 400), s.surface(320, 200)
	center, centerFlipped := s.surface(80, 16), s.surface(80, 16)
	s.part(center, logo, unionRegion(114, 0, 79, 15), 0, 0, 1, 1)
	s.part(centerFlipped, logo, unionRegion(114, 0, 79, 15), 0, 16, 1, -1)
	planes, err := scrolling.NewPlanes(scrolling.PlanesConfig{
		Slots: presets.TCBPlaneSlots(unionMultiplaneText, 32), Forms: presets.TCBScrollForms(), Visible: 30, PhaseStep: .02,
		Projection: unionMultiplaneProjection(),
	})
	if err != nil {
		s.err = err
		return
	}
	spec, _ := presets.FindFont("tcb-multi-plane-3d-scroller")
	metrics, err := spec.Build(font.Bounds())
	if err != nil {
		s.err = err
		return
	}
	renderer, err := scrolling.NewPlaneRenderer(scrolling.Face{Atlas: font, Metrics: metrics}, rasters)
	if err != nil {
		s.err = err
		return
	}
	s.closers = append(s.closers, renderer.Close)
	batch := composite.NewQuadBatch(64)
	bands, err := composite.NewBands(presets.UnionMountainBands())
	if err != nil {
		s.err = err
		return
	}
	profile := make([]float64, 40+804+810+160)
	for i := 40; i < 844; i++ {
		profile[i] = 8 * math.Sin(float64(i)*.05-2)
	}
	for i := 844; i < 1654; i++ {
		profile[i] = 8 * math.Sin(float64(i)*.15)
	}
	counter, next, rotation := 0, 0, 0.0
	s.render = func() {
		clearBlack(s.Canvas)
		background.Clear()
		stage.Clear()
		bands.Step()
		bands.DrawAt(background, mountains, 0, 0)
		s.draw(s.Canvas, background, 64, 60)
		counter++
		if counter > len(profile)-80 {
			counter = 0
		}
		batch.Begin(stage, logo)
		for i := 0; i < 32; i++ {
			batch.Rect(image.Rect(0, 16+i, 303, 17+i), float32(8+profile[counter+i]), float32(96+i), 303, 1)
		}
		batch.Flush()
		rotation += .08
		if rotation > 1 {
			rotation = -1
			next = 1 - next
		}
		img := center
		if next != 0 {
			img = centerFlipped
		}
		s.transform(stage, img, 160, 88, 1, rotation, 0, 40, 8, 1, ebiten.BlendSourceOver)
		if err := planes.Step(4); err != nil {
			s.err = err
			return
		}
		renderer.Draw(stage, planes.Points(), scrolling.PlaneDraw{ScaleX: 1, ScaleY: 1})
		s.transform(s.Canvas, stage, 64, 60, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
	}
}

func unionMultiplaneProjection() scrolling.PlaneProjection {
	// The composition positions glyph top-left corners. PlaneRenderer expects
	// their centers, so add half of the 32-by-33 font cell before projection.
	return scrolling.PlaneProjection{Focal: 250, Depth: 150, OriginX: -450, CenterX: 160, CenterY: 100, XBias: -16 + 32.0/2, YBias: -14 + 33.0/2, VerticalOffset: -4}
}

const unionMultiplaneText = " ^0                             " +
	"WOW, THIS DEMO SURE DOES LOOK GREAT..  BUT PERHAPS THE SCROLLINE LOOKS A BIT   TOO ORDINARY. " +
	"WELL, OKEY, LET US SWING IT UP AND DOWN. " +
	"^1 THIS IS THE LITTLE BIT OF EVERYTHING DEMO BY THE CAREBEARS. THERE ARE STAR RAY TYPE OF " +
	"BACKGROUND SCROLLERS, A DISTORTED TCB LOGO, " +
	"SOME GREAT MAD MAX MUSIC AND A SWINGING SCROLLINE OR..... PERHAPS EVEN MORE.............." +
	"^2...........  THIS IS BEGINNING TO LOOK " +
	"LIKE THE XXX INTERNATIONAL BALL DEMO SCREEN.                       " +
	"^3    BUT THEIR SCROLLINE WAS NOT THIS BIG. WE HOPE YOU DO NOT " +
	"THINK THAT WE HAVE TWO DIFFERENTLY SIZED FONTS. WE HAVE MANY MORE... ^4  " +
	"YEAH...  DO NOT LEAVE YET, THERE IS STILL MORE TO COME, JUST " +
	"WAIT AND SEE.  IF YOU THINK THIS IS HARD TO READ, WAIT TILL YOU HAVE " +
	"SEEN WHAT YOU ARE GOING TO SEE IN ABOUT THREE SECONDS.     " +
	"^5 THAT WAS NOT THREE SECONDS, BUT NOW YOU HAVE SEEN OUR THREE DIMENSIONAL " +
	"BENDING.. YOU MIGHT WONDER WHY WE HAVE NO PUNCTUATION EXCEPT " +
	"FOR THESE TWO ., . WE DO NOT EVEN HAVE THE LITTLE BLACK DOT BETWEEN HAVEN AND T, " +
	"HAVEN T, SEE... WELL, NOW THAT WE ARE OUT OF IDEAS WHAT " +
	"TO WRITE, WE CAN AS WELL EXPLAIN WHY. THE PROBLEM IS THAT ALL THE PART DEMOS " +
	"MUST WORK ON HALF A MEG AND EVERY CHARACTER TAKES ABOUT TEN " +
	"KILOBYTES. WE ARE GOING TO GREET SOME FOLKS NOW, SO LET US CHANGE WAVEFORM... " +
	"                        ^6             " +
	"MEGAGREETINGS GO TO ALL THE OTHER MEMBERS OF THE UNION. WE DO NOT FEEL " +
	"LIKE GREETING TO MUCH COZ WE DO NOT HAVE THOSE LITTLE BENT LINES, SO " +
	"WE CAN NOT MAKE COMMENTS. BUT JUST ONCE YOU WILL HAVE TO PRETEND YOU SAW " +
	"ONE OF THOSE, IT SHOULD HAVE COME INSTEAD OF THE SPACE BETWEEN " +
	"THE WORDS COOL AND YOUR. HERE WE GO... HELLO, AN COOL  YOUR NEW INTRO IS " +
	"REALLY SOMETHING .                    ^7 YOU WILL HAVE " +
	"TO READ IN THE MAIN SCROLLTEXT FOR MORE GREETINGS....  BYE.............. " +
	"                                             "
