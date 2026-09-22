package app

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/go-uniondemo/internal/screens"
)

// TourOptions controls the finite recording route. The complete introduction,
// loading credits and walking time are additional to these durations.
type TourOptions struct {
	ScreenDuration time.Duration
	MenuDuration   time.Duration
}

func DefaultTourOptions() TourOptions {
	return TourOptions{ScreenDuration: time.Minute, MenuDuration: 4 * time.Second}
}

// Tour drives the normal game with deterministic input and stops after visiting
// every door once. It never reads desktop input or changes the animation clock.
type Tour struct {
	*Game
	script tourScript
}

func NewTour(config Config, options TourOptions) (*Tour, error) {
	if options.ScreenDuration <= 0 || options.MenuDuration <= 0 {
		return nil, fmt.Errorf("tour: screen and menu durations must be positive")
	}
	config.Screen = "intro"
	config.Tour, config.SkipLoader, config.Touch = false, false, false
	g, err := New(config)
	if err != nil {
		return nil, err
	}
	return &Tour{Game: g, script: newTourScript(options, g.Rate())}, nil
}

func (t *Tour) Update() error {
	in, done, err := t.script.next(t.ScreenID(), t.hall.Selection())
	if err != nil {
		return err
	}
	if done {
		return ebiten.Termination
	}
	return t.Game.update(in)
}

// RecordingChapter names the actual rendered scene, including every return to
// the hall and each loading card. The video exporter records chapter changes.
func (t *Tour) RecordingChapter() string { return tourChapter(t.ScreenID()) }

func tourChapter(id string) string {
	if id == "menu" {
		return "The Union Hall"
	}
	loading := strings.HasPrefix(id, "loading:")
	d, ok := screens.Find(strings.TrimPrefix(id, "loading:"))
	if !ok {
		return id
	}
	if loading {
		return "Loading — " + d.Title
	}
	return d.Title
}

// tourScript only reasons about screen IDs and elapsed frames. Keeping it free
// of graphics and audio makes the entire route testable without rendering it.
type tourScript struct {
	phase                  string
	index, ticks, rate     int
	screenTicks, menuTicks int
	route                  []screens.Descriptor
}

func newTourScript(options TourOptions, rate int) tourScript {
	return tourScript{phase: "intro", rate: rate, screenTicks: durationTicks(options.ScreenDuration, rate), menuTicks: durationTicks(options.MenuDuration, rate), route: screens.Catalog()}
}

func durationTicks(duration time.Duration, rate int) int {
	return int(duration/time.Second)*rate + int((duration%time.Second*time.Duration(rate)+time.Second-1)/time.Second)
}

func (t *tourScript) next(id, selection string) (controls, bool, error) {
	unexpected := func() (controls, bool, error) {
		return controls{}, false, fmt.Errorf("tour: unexpected screen %q during %s", id, t.phase)
	}
	if t.phase == "intro" {
		if id == "intro" {
			t.ticks++
			if t.ticks > t.rate*120 {
				return controls{}, false, fmt.Errorf("tour: introduction did not finish")
			}
			return controls{}, false, nil
		}
		if id != "menu" {
			return unexpected()
		}
		t.phase, t.ticks = "menu", 1
	}
	if t.phase == "menu" {
		if id != "menu" {
			return unexpected()
		}
		// The first menu frame was drawn by the preceding return action.
		if t.ticks == 0 {
			t.ticks = 1
		}
		if t.ticks < t.menuTicks {
			t.ticks++
			return controls{}, false, nil
		}
		if t.index == len(t.route) {
			return controls{}, true, nil
		}
		t.phase, t.ticks = "walking", 0
	}
	target := t.route[t.index].ID
	if t.phase == "walking" {
		if id == "loading:"+target {
			t.phase, t.ticks = "loading", 1
		} else {
			if id != "menu" {
				return unexpected()
			}
			t.ticks++
			if t.ticks > 60*t.rate {
				return controls{}, false, fmt.Errorf("tour: walking to %s timed out", target)
			}
			if selection == target {
				return controls{enter: true}, false, nil
			}
			// Doors are ordered from left to right. Holding both directions also
			// brings Charly up to the doorway at the normal walking cadence.
			return controls{scene: screens.Input{Right: true, Up: true}}, false, nil
		}
	}
	if t.phase == "loading" {
		if id == "loading:"+target {
			t.ticks++
			if t.ticks > 30*t.rate {
				return controls{}, false, fmt.Errorf("tour: loading %s timed out", target)
			}
			return controls{}, false, nil
		}
		if id != target {
			return unexpected()
		}
		// The loading transition already drew the first frame of this screen.
		t.phase, t.ticks = "screen", 1
	}
	if id != target {
		return unexpected()
	}
	if t.ticks >= t.screenTicks {
		t.index++
		t.phase, t.ticks = "menu", 0
		return controls{back: true}, false, nil
	}
	in := tourScreenInput(target, t.ticks, t.screenTicks, t.rate)
	t.ticks++
	return controls{scene: in}, false, nil
}

func tourScreenInput(id string, tick, total, rate int) screens.Input {
	in := screens.Input{PointerX: Width / 2, PointerY: Height / 2}
	at := func(n, d int) bool { return tick == total*n/d }
	switch id {
	case "tnt3":
		// TNT is the initial object. Each selection uses the screen's existing
		// camera pullback and approach, leaving about nine seconds to rotate.
		for i, number := range []int{1, 3, 4, 5} {
			if at(i+1, 5) {
				in.Number = number
			}
		}
	case "hidden":
		seconds := float64(tick) / float64(rate)
		amplitude := math.Min(1, seconds/2)
		in.PointerX += amplitude * 220 * math.Sin(seconds*.8)
		in.PointerY += amplitude * 140 * math.Sin(seconds*1.1)
	case "starballs":
		for i, number := range []int{4, 7, 10} {
			if at(i+1, 4) {
				in.Number = number
			}
		}
	case "tnt2":
		if at(1, 4) {
			in.Right, in.Number = true, 1
		}
		if at(1, 3) {
			in.Right, in.Number = true, 2
		}
		if at(1, 2) {
			in.Action = true
		}
		if at(2, 3) {
			in.Left, in.Number = true, 4
		}
		if at(3, 4) {
			in.Left, in.Number = true, 3
		}
	case "diskcopier":
		// Keep the opening message, then allow the complete copy sequence to run.
		in.Action = tick == min(6*rate, max(1, total/10))
	}
	return in
}
