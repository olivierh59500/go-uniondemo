package screens

import (
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/effects"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

func init() { factories["tnt3"] = buildTNT3 }

type unionSolidFace struct {
	indices [4]int
	color   uint32
}

func unionFaces(faces []unionSolidFace) []effects.SolidFace {
	converted := make([]effects.SolidFace, len(faces))
	for i, face := range faces {
		converted[i] = effects.SolidFace{Indices: face.indices, Color: face.color}
	}
	return converted
}

func unionTNTObjects() ([]effects.SolidMeshModel, error) {
	sphere, err := effects.SolidSphereModel(effects.SolidSphereConfig{
		Radius: 80, Rows: 8, Columns: 16, PoleBand: true,
		BodyColors: [2]uint32{0xffffff, 0xff0000},
		PoleColors: [2]uint32{0x00ff00, 0x0000ff},
	})
	if err != nil {
		return nil, err
	}
	return []effects.SolidMeshModel{
		{Points: unionUnionPoints, Groups: [][]effects.SolidFace{unionFaces(unionUnionFaces)}},
		{Points: unionTntPoints, Groups: [][]effects.SolidFace{unionFaces(unionTntFaces)}},
		sphere,
		{Points: unionGliderPoints, Groups: [][]effects.SolidFace{unionFaces(unionGliderBaseFaces), unionFaces(unionGliderTopFaces)}},
		{Points: unionCarrierPoints, Groups: [][]effects.SolidFace{unionFaces(unionCarrierBottomFaces), unionFaces(unionCarrierPlaneFaces), unionFaces(unionCarrierTopFaces)}},
	}, nil
}

func buildTNT3(s *Scene) {
	stage, stars := s.surface(640, 400), s.image("stars.png")
	font := s.bitmap(s.image("fonts.png"), "union-tnt3")
	models, err := unionTNTObjects()
	if err != nil {
		s.err = err
		return
	}
	carousel, err := effects.NewSolidMeshCarousel(presets.UnionTNTMeshCarousel(models))
	if err != nil {
		s.err = err
		return
	}
	s.closers = append(s.closers, carousel.Close)
	captionConfig := presets.UnionTNTCaption(unionTNTText, font)
	caption, err := scrolling.New(scrolling.Config{Caption: &captionConfig})
	if err != nil {
		s.err = err
		return
	}
	s.closers = append(s.closers, caption.Close)
	pending, held := 1, false
	s.input = func(in Input) {
		pressed := in.Action || in.Left || in.Right
		selection := false
		if in.Number >= 1 && in.Number <= 5 {
			pending = in.Number - 1
			selection = true
		}
		if pressed && !held {
			if in.Left {
				pending = (pending + 4) % 5
			} else {
				pending = (pending + 1) % 5
			}
			selection = true
		}
		held = pressed
		if selection {
			s.err = carousel.Select(pending)
		}
	}
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		s.draw(stage, stars, 0, 0)
		if err := carousel.Step(); err != nil {
			s.err = err
			return
		}
		carousel.Draw(stage)
		if err := caption.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		caption.Draw(stage)
		s.draw(s.Canvas, stage, 64, 68)
	}
}

var unionTNTText = []string{
	"HEY GUYS !", "THE TNT-CREW IS VERY PROUD", "TO PRESENT YOU THEIR FAST", "3D-ROUTINES !!",
	"WE ARE NOT READY AS THERE", "ARE SOME ERRORS IN THE", "HIDDEN-FACE-ALGORITHM", "  ",
	"TRY THE KEYS 1..5", "FOR DIFFERENT OBJECTS", "  ", "1:  UNION-LOGO", "2:  TNT-LOGO", "3:  BALL", "4:  GLIDER", "5:  CARRIER", "  ",
	"PRESS ESCAPE TO EXIT", "  ", "PROGRAMED BY:", "JOJO AND HEXOGEN", "  ", "LOOK OUT FOR OUR", "3D-GAME !", "PERHAPS READY AT THE END OF 1989", "  ",
	"GREETINGS TO ALL GUYS OUT THERE !", "  ", "HAVE YOU ALREADY THE 7-TH", "TNT-DEMO ??", "AGAIN WITH NICE DIGI-SOUND.", "IT IS NOT SO BIG LIKE \"FNIL\"", "BUT VERY NICE!!",
	"THE LOUSY MUSIC IS FROM", "MAD MAX (TEX)", "DO YOU LIKE IT ?", "  ", "THE GREAT UNION-DEMO !!!", "  ",
}
