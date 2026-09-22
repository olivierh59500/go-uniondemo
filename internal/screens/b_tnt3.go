package screens

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/effects"
	"github.com/olivierh59500/democonstructionkit/geometry"
)

func init() { factories["tnt3"] = buildTNT3 }

type unionSolidFace struct {
	indices [4]int
	color   uint32
}

type unionSolidObject struct {
	points                  []geometry.Vec3
	groups                  [][]unionSolidFace
	rotation, camera, speed geometry.Vec3
}

func unionSolidMesh(points []geometry.Vec3, faces []unionSolidFace) effects.Mesh {
	mesh := effects.Mesh{Points: points}
	for _, face := range faces {
		c := face.color
		shade := color.NRGBA{R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c), A: 255}
		index := face.indices
		mesh.Triangles = append(mesh.Triangles, effects.Triangle{Indices: [3]int{index[0], index[1], index[2]}, Color: shade})
		if index[3] >= 0 {
			mesh.Triangles = append(mesh.Triangles, effects.Triangle{Indices: [3]int{index[0], index[2], index[3]}, Color: shade})
		}
	}
	return mesh
}

func unionTNTObjects() []unionSolidObject {
	spherePoints, sphereFaces := unionTNTSphere()
	standard := geometry.Vec3{X: .02, Y: .02, Z: .02}
	return []unionSolidObject{
		{points: unionUnionPoints, groups: [][]unionSolidFace{unionUnionFaces}, speed: standard},
		{points: unionTntPoints, groups: [][]unionSolidFace{unionTntFaces}, rotation: geometry.Vec3{Z: math.Pi}, camera: geometry.Vec3{Z: 100}, speed: standard},
		{points: spherePoints, groups: [][]unionSolidFace{sphereFaces}, speed: standard},
		{points: unionGliderPoints, groups: [][]unionSolidFace{unionGliderBaseFaces, unionGliderTopFaces}, rotation: geometry.Vec3{X: -math.Pi / 2}, speed: geometry.Vec3{X: .033, Y: .032, Z: .031}},
		{points: unionCarrierPoints, groups: [][]unionSolidFace{unionCarrierBottomFaces, unionCarrierPlaneFaces, unionCarrierTopFaces}, rotation: geometry.Vec3{X: -math.Pi / 2, Z: math.Pi / 3}, camera: geometry.Vec3{Y: 70, Z: 100}, speed: geometry.Vec3{Z: .02}},
	}
}

func buildTNT3(s *Scene) {
	stage, stars := s.surface(640, 400), s.image("stars.png")
	font := unionBitmap(s.image("fonts.png"), 16, 18)
	objects := unionTNTObjects()
	models := make([][]*effects.MeshEffect, len(objects))
	angles := objects[1].rotation
	var rx, ry, rz geometry.Rotation
	for i, object := range objects {
		for _, group := range object.groups {
			mesh, err := effects.NewMesh(unionSolidMesh(object.points, group), nil, geometry.Camera{Center: geometry.Vec2{X: 320, Y: 200}, Focal: 200 / math.Tan(25*math.Pi/360), Near: 1})
			if err != nil {
				s.err = err
				return
			}
			mesh.CullBackFaces = true
			// Preserve the object's local Euler order, then face the positive-Z camera.
			mesh.Deform = func(_ int, p geometry.Vec3, _ float64) geometry.Vec3 {
				p = rx.Apply(ry.Apply(rz.Apply(p)))
				return geometry.Vec3{X: p.X, Y: -p.Y, Z: -p.Z}
			}
			models[i] = append(models[i], mesh)
			s.closers = append(s.closers, mesh.Close)
		}
	}
	active, pending := 1, 1
	camera, distance := geometry.Vec3{Z: 10}, 0.0
	changing, held := false, false
	textIndex, textWait := 0, 200
	textY, textIncrement := -18.0, 2.0
	s.input = func(in Input) {
		pressed := in.Action || in.Left || in.Right
		if in.Number >= 1 && in.Number <= 5 {
			pending = in.Number - 1
			changing = true
		}
		if pressed && !held {
			if in.Left {
				pending = (pending + 4) % 5
			} else {
				pending = (pending + 1) % 5
			}
			changing = true
		}
		held = pressed
	}
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		rotating := distance >= 700
		if !rotating {
			distance += 10
		}
		if changing {
			if camera.Z < 10000 {
				camera.Z += 100
			} else {
				active = pending
				distance = 0
				camera = objects[active].camera
				angles = objects[active].rotation
				changing = false
			}
		}
		s.draw(stage, stars, 0, 0)
		rx = geometry.RotateXYZ(geometry.Vec3{X: angles.X})
		ry = geometry.RotateXYZ(geometry.Vec3{Y: angles.Y})
		rz = geometry.RotateXYZ(geometry.Vec3{Z: angles.Z})
		for _, mesh := range models[active] {
			mesh.Transform.Position = geometry.Vec3{X: -camera.X, Y: camera.Y, Z: camera.Z + distance}
			if err := mesh.Update(kit.Frame{}); err != nil {
				s.err = err
				return
			}
			mesh.Draw(stage)
		}
		if rotating {
			angles = angles.Add(objects[active].speed)
		}
		stage.SubImage(image.Rect(0, 0, 640, 18)).(*ebiten.Image).Fill(color.Black)
		line := unionTNTText[textIndex]
		font.Print(stage, line, 320-float64(len(line)*8), textY, 1, 1)
		if textWait == 100 {
			textY += textIncrement
		}
		if textY >= 0 {
			textY = 0
			textWait--
			if textWait <= 0 {
				textIncrement = -2
				textWait = 100
			}
		}
		if textY <= -18 {
			textY = -18
			textWait--
			if textWait <= 0 {
				textIncrement = 2
				textWait = 100
				textIndex = (textIndex + 1) % len(unionTNTText)
			}
		}
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
