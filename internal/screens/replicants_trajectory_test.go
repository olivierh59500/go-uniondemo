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
	if first.X != 257 || first.Y != 237 || last.X != 517 || last.Y != 237 {
		t.Fatalf("unexpected resting sprite row: first=%+v last=%+v", first, last)
	}
	for _, seconds := range []float64{.4, 1.2, 2.4, 3, 6.4, 8, 10, 12, 14, 16, 18, 20, 22, 24} {
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
		if end, start := formation.At(26-1.0/60, index), formation.At(26, index); math.Hypot(end.X-start.X, end.Y-start.Y) > 1 {
			t.Fatalf("sprite %d jumps when the choreography loops: %+v to %+v", index, end, start)
		}
	}
	for frame := 0; frame < 26*60; frame++ {
		for index := 0; index < 14; index++ {
			before := formation.At(float64(frame)/60, index)
			after := formation.At(float64(frame+1)/60, index)
			// The opening sweep is fast but every key transition is continuous.
			if distance := math.Hypot(after.X-before.X, after.Y-before.Y); distance > 60 {
				t.Fatalf("sprite %d jumps %.1f pixels at frame %d", index, distance, frame)
			}
		}
	}
}

func TestReplicantsEntryCompressesAndSwingsAcrossBothPanels(t *testing.T) {
	formation, err := replicantsFormation()
	if err != nil {
		t.Fatal(err)
	}
	check := func(seconds, wantCenter, wantWidth, wantY float64) {
		t.Helper()
		left, right := formation.At(seconds, 0), formation.At(seconds, 13)
		center, width := (left.X+right.X)/2, right.X-left.X
		if math.Abs(center-wantCenter) > 1 || math.Abs(width-wantWidth) > 1 || math.Abs(left.Y-wantY) > 1 {
			t.Fatalf("entry pose at %.1fs center/width/Y = %.2f/%.2f/%.2f, want %.2f/%.2f/%.2f", seconds, center, width, left.Y, wantCenter, wantWidth, wantY)
		}
	}
	check(0, 387, 260, 237)
	check(.6, 263.5, 273, 250)
	check(1, 489.25, 318.5, 250)
	check(1.4, 603.55, -120.9, 238)
	check(2, 303.75, -318.5, 237)
	check(2.2, 217.25, -175.5, 240)
	check(2.4, 189.45, 120.9, 250)
	check(3, 489.25, 318.5, 250)
	check(3.4, 603.55, -120.9, 238)
	check(4, 303.75, -318.5, 237)
	check(4.4, 189.45, 120.9, 250)
	check(5, 489.25, 318.5, 250)
	check(5.4, 603.55, -120.9, 238)
	check(6, 303.75, -318.5, 237)
	check(6.4, 189.45, 120.9, 250)
	check(7, 379.75, 435.5, 232)
	check(8, 415.25, 370.5, 184)
	check(10, 374.75, 435.5, 231)
	check(11, 343.75, 409.5, 190)
	check(18, 390.5, 221, 238)
	check(20, 335.05, -341.9, 240)
	check(25, 383.75, 435.5, 237)
	if middle := formation.At(8, 6); math.Abs(middle.Y-151.4) > 1 {
		t.Fatalf("the authored arch was not applied: %+v", middle)
	}
	for index := 0; index < 14; index++ {
		if first, repeat := formation.At(1.4, index), formation.At(3.4, index); first != repeat {
			t.Fatalf("two-second sweep changed at sprite %d: %+v to %+v", index, first, repeat)
		}
	}
}
