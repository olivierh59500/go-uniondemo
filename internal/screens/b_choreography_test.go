package screens

import (
	"math"
	"testing"

	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/timeline"
)

func TestUnionSolidObjects(t *testing.T) {
	objects := unionTNTObjects()
	wantPoints := []int{28, 48, 512, 7, 29}
	wantFaces := []int{45, 40, 112, 10, 20}
	if len(objects) != 5 {
		t.Fatalf("objects = %d, want 5", len(objects))
	}
	for i, object := range objects {
		if len(object.points) != wantPoints[i] {
			t.Errorf("object %d: points = %d, want %d", i, len(object.points), wantPoints[i])
		}
		faceCount := 0
		for _, group := range object.groups {
			faceCount += len(group)
			mesh := unionSolidMesh(object.points, group)
			for _, triangle := range mesh.Triangles {
				for _, index := range triangle.Indices {
					if index < 0 || index >= len(object.points) {
						t.Errorf("object %d: invalid vertex index %d", i, index)
					}
				}
				if triangle.Color.A != 255 {
					t.Errorf("object %d: unexpected transparent material", i)
				}
			}
		}
		if faceCount != wantFaces[i] {
			t.Errorf("object %d: faces = %d, want %d", i, faceCount, wantFaces[i])
		}
		for _, p := range object.points {
			if math.IsNaN(p.X+p.Y+p.Z) || math.IsInf(p.X+p.Y+p.Z, 0) {
				t.Errorf("object %d: nonfinite vertex", i)
			}
		}
	}
}

func TestUnionMultiplaneTopLeftProjection(t *testing.T) {
	plane, err := scrolling.NewPlanes(scrolling.PlanesConfig{
		Slots: []scrolling.PlaneSlot{{Rune: 'A', Advance: 32, Form: 0}},
		Forms: presets.TCBScrollForms(), Projection: unionMultiplaneProjection(), Visible: 1, PhaseStep: .02,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := plane.Step(4); err != nil {
		t.Fatal(err)
	}
	point := plane.Points()[0]
	// Check the visible atlas corner, including the renderer's centered anchor.
	wantX := (-450.0-16)*250/400 + 160
	wantY := (55*math.Cos(1.5)-4-14)*250/400 + 100
	if x, y := point.X-16*point.Scale, point.Y-16.5*point.Scale; math.Abs(x-wantX) > 1e-10 || math.Abs(y-wantY) > 1e-10 {
		t.Fatalf("first glyph top-left = (%g, %g), want (%g, %g)", x, y, wantX, wantY)
	}
}

func TestUnionCopyStages(t *testing.T) {
	ranges, err := timeline.NewCueRanges(presets.UnionDiskCopierCueRanges())
	if err != nil {
		t.Fatal(err)
	}
	palette, err := timeline.NewSteppedEnvelope(presets.UnionDiskCopierFade())
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		time      float64
		operation int
		message   bool
	}{
		{0, -1, true}, {99.5, -1, true}, {100, -1, false}, {120, -1, false}, {120.5, 0, true}, {359.5, 0, true},
		{360, -1, false}, {380.5, -1, true}, {480, -1, false}, {500.5, 1, true}, {739.5, 1, true},
		{740, -1, false}, {760.5, 2, true}, {999.5, 2, true}, {1000, -1, false}, {1000.5, -1, true},
		{1100, -1, false}, {1120, -1, false}, {1120.5, -1, true}, {1220.5, -1, true},
	} {
		line, operation := "", -1
		if index, _, active := ranges.At(test.time); active {
			line, operation = unionCopyStages[index].line, unionCopyStages[index].operation
		}
		if operation != test.operation || (line != "") != test.message {
			t.Errorf("time %g: operation %d, message %q", test.time, operation, line)
		}
	}
	for _, test := range []struct {
		time  float64
		shade int
	}{{0, 7}, {2, 6}, {12, 1}, {14, 0}, {84, 0}, {86, 1}, {96, 6}, {98, 7}, {100, 7}} {
		if got := palette.At(test.time, 0, 100); got != test.shade {
			t.Errorf("fade at %g = %d, want %d", test.time, got, test.shade)
		}
	}
}
