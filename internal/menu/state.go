// Package menu presents the Union's walkable screen selection hall.
package menu

// Input contains the held movement controls and the enter action for one tick.
type Input struct {
	Left, Right, Up, Down, Enter bool
}

type door struct {
	id       string
	min, max int
}

var doors = [...]door{
	{"beatdis", 340, 392},
	{"delta", 980, 1042},
	{"tnt3", 1264, 1328},
	{"wow", 1488, 1548},
	{"hidden", 1820, 1868},
	{"starballs", 2232, 2288},
	{"replicants", 2640, 2708},
	{"tnt2", 3152, 3208},
	{"level16", 3568, 3628},
	{"multiplane", 3852, 3916},
	{"diskcopier", 5008, 5072},
}

const hallLength = 5506

// state contains only deterministic animation and navigation data. Rendering
// never advances it, even when the display refresh rate differs from 60 Hz.
type state struct {
	tick                               uint64
	backX, bannerX, charleyY           int
	charleyFrame, charleyDelay, facing int
	moveDelay                          int
	panoramaX, coverWidth              int
	logoScale, logoStep                float64
	logoWait                           int
	colorIndex, colorDelay             int
}

func newState() state {
	return state{
		charleyY: 90, facing: 1,
		coverWidth: 840,
		logoScale:  1, logoStep: -0.04, logoWait: 1000,
	}
}

func (s *state) update(in Input) string {
	s.tick++
	s.panoramaX -= 2
	if s.panoramaX <= -540 {
		s.panoramaX = -28
	}
	s.coverWidth = max(0, s.coverWidth-5)
	if s.colorDelay >= 3 {
		s.colorIndex = (s.colorIndex + 1) % len(hallColors)
		s.colorDelay = 0
	}
	s.colorDelay++
	if s.logoWait <= 5 {
		s.logoScale += s.logoStep
		if s.logoScale <= 0 {
			s.logoScale = 0
			s.logoStep = 0.04
		}
	}
	if s.logoWait <= 0 && s.logoScale >= 1 {
		s.logoScale = 1
		s.logoStep = -0.04
		s.logoWait = 1000
	}
	s.logoWait--

	dx, dy := 0, 0
	if in.Left != in.Right {
		if in.Left {
			dx = -1
		} else {
			dx = 1
		}
	}
	if in.Up != in.Down {
		if in.Up {
			dy = -1
		} else {
			dy = 1
		}
	}
	if dx == 0 && dy == 0 {
		s.moveDelay = 0
	} else if s.moveDelay == 0 {
		// A movement action every two ticks retains the walking cadence while
		// using continuous keyboard, gamepad or touch input.
		s.move(dx, dy)
		s.moveDelay = 1
	} else {
		s.moveDelay--
	}
	if in.Enter {
		return s.selection()
	}
	return ""
}

func (s *state) move(dx, dy int) {
	if dx != 0 {
		s.facing = dx
		s.backX -= dx * 5
		if s.backX <= -hallLength {
			s.backX = 0
		} else if s.backX >= 0 {
			s.backX = -hallLength
		}
		s.bannerX -= dx * 3
		if dx > 0 && s.bannerX <= -94 {
			s.bannerX = -78
		} else if dx < 0 && s.bannerX >= -78 {
			s.bannerX = -94
		}
		s.animateCharacter()
	}
	if dy != 0 {
		s.charleyY = min(120, max(60, s.charleyY+dy*4))
		s.animateCharacter()
	}
}

func (s *state) animateCharacter() {
	s.charleyDelay++
	if s.charleyDelay >= 5 {
		s.charleyFrame = (s.charleyFrame + 1) % 8
		s.charleyDelay = 0
	}
}

func (s *state) selection() string {
	if s.charleyY > 62 {
		return ""
	}
	x := -s.backX
	for _, d := range doors {
		if x >= d.min && x <= d.max {
			return d.id
		}
	}
	return ""
}

func (s *state) positionDoor(id string) {
	for _, d := range doors {
		if d.id == id {
			s.backX = -(d.min + d.max) / 2
			s.charleyY = 60
			s.moveDelay = 0
			return
		}
	}
}
