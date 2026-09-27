package screens

import (
	"math"
	"testing"
)

func TestReplicantsKeepsFastSweepsAndIndividualArcPaths(t *testing.T) {
	formation, err := replicantsFormation()
	if err != nil {
		t.Fatal(err)
	}
	if first, last := formation.At(0, 0), formation.At(0, 13); first.X != 257 || first.Y != 237 || last.X != 517 || last.Y != 237 {
		t.Fatalf("unexpected opening stack: %+v to %+v", first, last)
	}
	if formation.At(1.4, 0).X <= formation.At(1.4, 13).X || formation.At(2.4, 0).X >= formation.At(2.4, 13).X {
		t.Fatal("opening stacks no longer reverse from right to left")
	}
	firstStart, firstEnd := formation.At(15.2, 0), formation.At(15.6, 0)
	lastStart, lastEnd := formation.At(15.2, 13), formation.At(15.6, 13)
	if firstEnd.X-firstStart.X < 200 || firstStart.Y-firstEnd.Y < 100 || lastEnd.Y-lastStart.Y < 100 {
		t.Fatalf("letters lost their rapid crossing arcs: first %+v to %+v; last %+v to %+v", firstStart, firstEnd, lastStart, lastEnd)
	}
	for second := 7; second < 25; second++ {
		maximum := 0.0
		for index := 0; index < 14; index++ {
			a, b := formation.At(float64(second), index), formation.At(float64(second+1), index)
			maximum = math.Max(maximum, math.Hypot(a.X-b.X, a.Y-b.Y))
		}
		if maximum < 20 {
			t.Fatalf("the letter motion nearly stopped around %ds: %.1f pixels", second, maximum)
		}
	}
}

func TestReplicantsFormationStaysVisibleAndJoinsItsLoop(t *testing.T) {
	formation, err := replicantsFormation()
	if err != nil {
		t.Fatal(err)
	}
	for index := 0; index < 14; index++ {
		if formation.At(25, index) != formation.At(0, index) {
			t.Fatalf("sprite %d jumps at the loop boundary", index)
		}
	}
	for frame := 0; frame < 25*60; frame++ {
		before, after := float64(frame)/60, float64(frame+1)/60
		for index := 0; index < 14; index++ {
			a, b := formation.At(before, index), formation.At(after, index)
			if a.X < 0 || a.X > 736 || a.Y < 120 || a.Y > 355 {
				t.Fatalf("sprite %d leaves the visible band at frame %d: %+v", index, frame, a)
			}
			if distance := math.Hypot(b.X-a.X, b.Y-a.Y); distance > 16 {
				t.Fatalf("sprite %d jumps %.1f pixels at frame %d", index, distance, frame)
			}
		}
	}
}
