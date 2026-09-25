package screens

import (
	"math"
	"testing"
)

func TestReplicantsFormationCoversReferenceMotionWithoutLoopJump(t *testing.T) {
	formation, err := replicantsFormation()
	if err != nil {
		t.Fatal(err)
	}
	first, last := formation.At(0, 0), formation.At(0, 13)
	if first.X != 162 || first.Y != 230 || last.X != 571.5 || last.Y != 230 {
		t.Fatalf("unexpected resting sprite row: first=%+v last=%+v", first, last)
	}
	for _, seconds := range []float64{.8, 3, 6, 10, 13.5, 17, 21, 23} {
		largest := 0.0
		for index := 0; index < 14; index++ {
			pose := formation.At(seconds, index)
			rest := formation.At(0, index)
			largest = math.Max(largest, math.Hypot(pose.X-rest.X, pose.Y-rest.Y))
			if pose.X < 0 || pose.X > 736 || pose.Y < 120 || pose.Y > 355 {
				t.Fatalf("sprite %d leaves useful screen space at %.1fs: %+v", index, seconds, pose)
			}
		}
		if largest < 20 {
			t.Fatalf("trajectory at %.1fs became nearly static: max displacement %.1f", seconds, largest)
		}
	}
	for index := 0; index < 14; index++ {
		if end, start := formation.At(25-1.0/60, index), formation.At(25, index); math.Hypot(end.X-start.X, end.Y-start.Y) > 1 {
			t.Fatalf("sprite %d jumps when the choreography loops: %+v to %+v", index, end, start)
		}
	}
	for frame := 0; frame < 25*60; frame++ {
		for index := 0; index < 14; index++ {
			before := formation.At(float64(frame)/60, index)
			after := formation.At(float64(frame+1)/60, index)
			if distance := math.Hypot(after.X-before.X, after.Y-before.Y); distance > 20 {
				t.Fatalf("sprite %d jumps %.1f pixels at frame %d", index, distance, frame)
			}
		}
	}
}
