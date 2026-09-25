// Package menu presents the Union's walkable screen selection hall.
package menu

import (
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/timeline"
)

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
	tick                     uint64
	backX, bannerX, charleyY int
	charleyFrame, facing     int
	moveDelay                int
	panoramaX, coverWidth    int
	logoScale                float64
	colorIndex               int
	panorama                 *motion.WrapBank
	cover                    motion.LinearTick
	paletteCycle             *timeline.PacedIndex
	characterCycle           *timeline.PacedIndex
	logoMotion               *motion.HoldBounce
	walkMotion               *motion.WalkParallax
}

func newState() state {
	panorama, err := motion.NewWrapBank(presets.UnionMenuPanoramaWrap())
	if err != nil {
		panic(err)
	}
	cover, err := motion.NewLinearTick(presets.UnionMenuCoverWidth())
	if err != nil {
		panic(err)
	}
	paletteCycle, err := timeline.NewPacedIndex(presets.UnionMenuPaletteCycle(len(hallColors)))
	if err != nil {
		panic(err)
	}
	characterCycle, err := timeline.NewPacedIndex(presets.UnionMenuCharacterCycle())
	if err != nil {
		panic(err)
	}
	logoMotion, err := motion.NewHoldBounce(presets.UnionMenuLogoBounce())
	if err != nil {
		panic(err)
	}
	walkMotion, err := motion.NewWalkParallax(presets.UnionMenuWalkParallax(hallLength))
	if err != nil {
		panic(err)
	}
	return state{
		charleyY: 90, facing: 1,
		coverWidth: cover.At(0), logoScale: logoMotion.At(),
		panorama: panorama, cover: cover, paletteCycle: paletteCycle,
		characterCycle: characterCycle, logoMotion: logoMotion,
		walkMotion: walkMotion,
	}
}

func (s *state) update(in Input) string {
	s.tick++
	s.panorama.Step()
	s.panoramaX = int(s.panorama.At(0))
	s.coverWidth = s.cover.At(int(s.tick))
	s.colorIndex = s.paletteCycle.Step()
	s.logoScale = s.logoMotion.Step()

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
		s.walkMotion.Advance(dx)
		s.backX = int(s.walkMotion.At(0))
		s.bannerX = int(s.walkMotion.At(1))
		s.animateCharacter()
	}
	if dy != 0 {
		s.charleyY = min(120, max(60, s.charleyY+dy*4))
		s.animateCharacter()
	}
}

func (s *state) animateCharacter() {
	s.charleyFrame = s.characterCycle.Step()
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
			s.walkMotion.Set(0, float64(s.backX))
			s.charleyY = 60
			s.moveDelay = 0
			return
		}
	}
}
