package screens

import (
	"math"

	"github.com/olivierh59500/democonstructionkit/effects"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

func init() { factories["tnt3"] = buildTNT3 }

type unionSolidFace struct {
	indices [4]int
	color   uint32
}

type unionSolidObject struct {
	points []geometry.Vec3
	groups [][]unionSolidFace
}

func unionSolidMesh(points []geometry.Vec3, faces []unionSolidFace) effects.Mesh {
	return effects.SolidMesh(points, unionFaces(faces))
}

func unionFaces(faces []unionSolidFace) []effects.SolidFace {
	converted := make([]effects.SolidFace, len(faces))
	for i, face := range faces {
		converted[i] = effects.SolidFace{Indices: face.indices, Color: face.color}
	}
	return converted
}

func unionTNTObjects() []unionSolidObject {
	spherePoints, sphereFaces := unionTNTSphere()
	return []unionSolidObject{
		{points: unionUnionPoints, groups: [][]unionSolidFace{unionUnionFaces}},
		{points: unionTntPoints, groups: [][]unionSolidFace{unionTntFaces}},
		{points: spherePoints, groups: [][]unionSolidFace{sphereFaces}},
		{points: unionGliderPoints, groups: [][]unionSolidFace{unionGliderBaseFaces, unionGliderTopFaces}},
		{points: unionCarrierPoints, groups: [][]unionSolidFace{unionCarrierBottomFaces, unionCarrierPlaneFaces, unionCarrierTopFaces}},
	}
}

func buildTNT3(s *Scene) {
	stage, stars := s.surface(640, 400), s.image("stars.png")
	font := s.bitmap(s.image("fonts.png"), "union-tnt3")
	objects := unionTNTObjects()
	models := make([]effects.SolidMeshModel, len(objects))
	for i, object := range objects {
		models[i].Points = object.points
		models[i].Groups = make([][]effects.SolidFace, len(object.groups))
		for j, group := range object.groups {
			models[i].Groups[j] = unionFaces(group)
		}
	}
	carousel, err := effects.NewSolidMeshCarousel(presets.UnionTNTMeshCarousel(models))
	if err != nil {
		s.err = err
		return
	}
	s.closers = append(s.closers, carousel.Close)
	caption, err := scrolling.NewCaptionCarousel(presets.UnionTNTCaption(unionTNTText, font))
	if err != nil {
		s.err = err
		return
	}
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
		caption.Draw(stage)
		caption.Step()
		s.draw(s.Canvas, stage, 64, 68)
	}
}

func unionTNTSphere() ([]geometry.Vec3, []unionSolidFace) {
	points := make([]geometry.Vec3, 0, 512)
	faces := make([]unionSolidFace, 0, 112)
	point := func(theta, phi float64) geometry.Vec3 {
		theta *= math.Pi / 180
		phi *= math.Pi / 180
		return geometry.Vec3{X: 80 * math.Cos(theta) * math.Cos(phi), Y: 80 * math.Cos(theta) * math.Sin(phi), Z: -80 * math.Sin(theta)}
	}
	for row := 0; row < 8; row++ {
		theta := -90 + float64(row)*22.5
		for column := 0; column < 16; column++ {
			phi := float64(column) * 22.5
			start := len(points)
			a, b, c, d := point(theta, phi), point(theta+22.5, phi), point(theta+22.5, phi+22.5), point(theta, phi+22.5)
			if row == 0 {
				d = c
			}
			points = append(points, a, b, c, d)
			shade := uint32(0xffffff)
			if (column+row)%2 != 0 {
				shade = 0xff0000
			}
			if row > 0 && row < 7 {
				faces = append(faces, unionSolidFace{[4]int{start, start + 1, start + 2, start + 3}, shade})
			}
		}
	}
	for i := 0; i < 16; i++ {
		a, b, c, d := i*4+1, i*4+2, 452+i*4, 448+i*4
		if i == 15 {
			b, c, d = 1, 448, 508
		}
		shade := uint32(0x00ff00)
		if i%2 != 0 {
			shade = 0x0000ff
		}
		faces = append(faces, unionSolidFace{[4]int{a, b, c, d}, shade})
	}
	return points, faces
}

var unionTNTText = []string{
	"HEY GUYS !", "THE TNT-CREW IS VERY PROUD", "TO PRESENT YOU THEIR FAST", "3D-ROUTINES !!",
	"WE ARE NOT READY AS THERE", "ARE SOME ERRORS IN THE", "HIDDEN-FACE-ALGORITHM", "  ",
	"TRY THE KEYS 1..5", "FOR DIFFERENT OBJECTS", "  ", "1:  UNION-LOGO", "2:  TNT-LOGO", "3:  BALL", "4:  GLIDER", "5:  CARRIER", "  ",
	"PRESS ESCAPE TO EXIT", "  ", "PROGRAMED BY:", "JOJO AND HEXOGEN", "  ", "LOOK OUT FOR OUR", "3D-GAME !", "PERHAPS READY AT THE END OF 1989", "  ",
	"GREETINGS TO ALL GUYS OUT THERE !", "  ", "HAVE YOU ALREADY THE 7-TH", "TNT-DEMO ??", "AGAIN WITH NICE DIGI-SOUND.", "IT IS NOT SO BIG LIKE \"FNIL\"", "BUT VERY NICE!!",
	"THE LOUSY MUSIC IS FROM", "MAD MAX (TEX)", "DO YOU LIKE IT ?", "  ", "THE GREAT UNION-DEMO !!!", "  ",
}
