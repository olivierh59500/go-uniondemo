package menu

import "testing"

func TestDoorSelectionAndEntry(t *testing.T) {
	s := newState()
	if got := s.update(Input{Enter: true}); got != "" {
		t.Fatalf("initial position selected %q", got)
	}
	for _, d := range doors {
		t.Run(d.id, func(t *testing.T) {
			s.positionDoor(d.id)
			if got := s.selection(); got != d.id {
				t.Fatalf("positioned at %q, selected %q", d.id, got)
			}
			if got := s.update(Input{}); got != "" {
				t.Fatalf("entered %q without an enter action", got)
			}
			if got := s.update(Input{Enter: true}); got != d.id {
				t.Fatalf("enter selected %q, want %q", got, d.id)
			}
			for _, x := range []int{d.min, d.max} {
				s.backX = -x
				s.charleyY = 62
				if got := s.selection(); got != d.id {
					t.Fatalf("door boundary x=%d selected %q", x, got)
				}
			}
			for _, x := range []int{d.min - 1, d.max + 1} {
				s.backX = -x
				if got := s.selection(); got != "" {
					t.Fatalf("outside door x=%d selected %q", x, got)
				}
			}
			s.positionDoor(d.id)
			s.charleyY = 63
			if got := s.update(Input{Enter: true}); got != "" {
				t.Fatalf("standing in front of the door selected %q", got)
			}
		})
	}
	before := s
	s.positionDoor("unknown")
	if s != before {
		t.Fatal("unknown door changed menu state")
	}
}

func TestWalkingCadenceAndCharacterAnimation(t *testing.T) {
	s := newState()
	for i := 0; i < 60; i++ {
		s.update(Input{Right: true})
	}
	if s.backX != -150 || s.charleyFrame != 6 || s.facing != 1 {
		t.Fatalf("one second walk: x=%d frame=%d facing=%d", s.backX, s.charleyFrame, s.facing)
	}
	s.update(Input{})
	s.update(Input{Left: true})
	if s.backX != -145 || s.facing != -1 {
		t.Fatalf("new direction should respond immediately: x=%d facing=%d", s.backX, s.facing)
	}
	x, y, frame := s.backX, s.charleyY, s.charleyFrame
	for i := 0; i < 60; i++ {
		s.update(Input{Left: true, Right: true, Up: true, Down: true})
	}
	if s.backX != x || s.charleyY != y || s.charleyFrame != frame {
		t.Fatal("opposite held controls moved or animated the character")
	}
}

func TestHallWrapAndVerticalLimits(t *testing.T) {
	s := newState()
	s.update(Input{Left: true})
	if s.backX != -hallLength {
		t.Fatalf("left wrap: %d", s.backX)
	}
	s.update(Input{})
	s.update(Input{Right: true})
	if s.backX != 0 {
		t.Fatalf("right wrap: %d", s.backX)
	}
	for i := 0; i < 100; i++ {
		s.update(Input{Up: true})
	}
	if s.charleyY != 60 {
		t.Fatalf("upper limit: %d", s.charleyY)
	}
	for i := 0; i < 100; i++ {
		s.update(Input{Down: true})
	}
	if s.charleyY != 120 {
		t.Fatalf("lower limit: %d", s.charleyY)
	}
}

func TestAnimationBoundsAndRepeat(t *testing.T) {
	s := newState()
	compressed := false
	expanded := false
	for i := 0; i < 5000; i++ {
		s.update(Input{})
		if s.logoScale < 0 || s.logoScale > 1 {
			t.Fatalf("logo scale %f at tick %d", s.logoScale, i)
		}
		if s.logoScale == 0 {
			compressed = true
		}
		if compressed && s.logoScale == 1 {
			expanded = true
		}
		if s.panoramaX <= -540 || s.panoramaX > 0 || s.coverWidth < 0 {
			t.Fatalf("scroll bounds at tick %d: panorama=%d cover=%d", i, s.panoramaX, s.coverWidth)
		}
		if s.colorIndex < 0 || s.colorIndex >= len(hallColors) {
			t.Fatalf("palette index %d", s.colorIndex)
		}
	}
	if !compressed || !expanded || s.coverWidth != 0 {
		t.Fatalf("animation did not complete a cycle: compressed=%v expanded=%v cover=%d", compressed, expanded, s.coverWidth)
	}
}
