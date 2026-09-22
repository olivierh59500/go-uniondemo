package app

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/olivierh59500/go-uniondemo/assets"
	"github.com/olivierh59500/go-uniondemo/internal/loader"
	"github.com/olivierh59500/go-uniondemo/internal/menu"
	"github.com/olivierh59500/go-uniondemo/internal/screens"
)

func TestTourScriptCompleteTimeline(t *testing.T) {
	for _, rate := range []int{50, 60} {
		options := DefaultTourOptions()
		script := newTourScript(options, rate)
		id, selection := "intro", ""
		counts := map[string]int{}
		var chapters []string
		menuWait, walk := 0, 0
		finished := false
		for frame := 0; frame < 30*60*rate; frame++ {
			in, done, err := script.next(id, selection)
			if err != nil {
				t.Fatal(err)
			}
			if done {
				finished = true
				break
			}
			switch {
			case id == "intro":
				if in.enter || in.back || in.scene.Action {
					t.Fatal("tour skipped the introduction")
				}
				if counts[id] == durationTicks(screens.IntroDuration, rate) {
					id, menuWait = "menu", 1
				}
			case id == "menu":
				if in.enter {
					if selection != script.route[script.index].ID {
						t.Fatal("entered the wrong door")
					}
					id, walk = "loading:"+selection, 0
				} else if in.scene.Right {
					if menuWait != durationTicks(options.MenuDuration, rate) {
						t.Fatalf("menu rested for %d frames", menuWait)
					}
					if !in.scene.Up {
						t.Fatal("walk did not approach the doorway")
					}
					walk++
					// Wrong doors along the path must be passed without entry.
					selection = "another-door"
					if walk == 17 {
						selection = script.route[script.index].ID
					}
				} else {
					menuWait++
				}
			case strings.HasPrefix(id, "loading:"):
				if in.enter || in.back {
					t.Fatal("tour skipped the loading credits")
				}
				if counts[id] == durationTicks(loader.Duration, rate) {
					id = strings.TrimPrefix(id, "loading:")
				}
			default:
				if in.back {
					if counts[id] != durationTicks(options.ScreenDuration, rate) {
						t.Fatalf("%s lasted %d frames", id, counts[id])
					}
					id, menuWait = "menu", 1
				}
			}
			counts[id]++
			name := tourChapter(id)
			if len(chapters) == 0 || chapters[len(chapters)-1] != name {
				chapters = append(chapters, name)
			}
		}
		if !finished || id != "menu" || menuWait != durationTicks(options.MenuDuration, rate) {
			t.Fatalf("tour failed to finish with a complete final menu: finished=%v screen=%s wait=%d", finished, id, menuWait)
		}
		wantChapters := []string{tourChapter("intro"), tourChapter("menu")}
		for _, d := range screens.Catalog() {
			wantChapters = append(wantChapters, tourChapter("loading:"+d.ID), tourChapter(d.ID), tourChapter("menu"))
		}
		if !reflect.DeepEqual(chapters, wantChapters) {
			t.Fatalf("chapter route differs:\n got: %v\nwant: %v", chapters, wantChapters)
		}
	}
}

func TestTourWalksToEveryRealDoor(t *testing.T) {
	// Only the hall's state advances. No screen rendering commands accumulate.
	hall, err := menu.New(assets.Files)
	if err != nil {
		t.Fatal(err)
	}
	defer hall.Close()
	script := newTourScript(TourOptions{ScreenDuration: time.Second, MenuDuration: time.Second}, 60)
	id, walkingFrames := "menu", 0
	for tick := 0; tick < 6000; tick++ {
		selection := hall.Selection()
		in, done, err := script.next(id, selection)
		if err != nil {
			t.Fatal(err)
		}
		if done {
			if walkingFrames < 1900 || walkingFrames > 2100 {
				t.Fatalf("unexpected complete hall walking time: %d frames", walkingFrames)
			}
			return
		}
		switch {
		case id == "menu":
			entered := hall.Update(menu.Input{Left: in.scene.Left, Right: in.scene.Right, Up: in.scene.Up, Down: in.scene.Down, Enter: in.enter})
			if in.scene.Right {
				walkingFrames++
			}
			if entered != "" {
				id = "loading:" + entered
			}
		case strings.HasPrefix(id, "loading:"):
			if script.phase == "loading" {
				id = strings.TrimPrefix(id, "loading:")
			}
		case in.back:
			if hall.Selection() != id {
				t.Fatalf("return from %s lost Charly's position", id)
			}
			id = "menu"
		}
	}
	t.Fatal("tour did not finish walking all eleven real doors")
}

func TestTourDemonstratesInteractiveScreens(t *testing.T) {
	const rate, total = 60, 3600
	var objects, counts, layerKeys []int
	copyActions := 0
	previous := tourScreenInput("hidden", 0, total, rate)
	for tick := 1; tick < total; tick++ {
		for _, test := range []struct {
			id  string
			dst *[]int
		}{{"tnt3", &objects}, {"starballs", &counts}, {"tnt2", &layerKeys}} {
			if number := tourScreenInput(test.id, tick, total, rate).Number; number != 0 {
				*test.dst = append(*test.dst, number)
			}
		}
		if tourScreenInput("diskcopier", tick, total, rate).Action {
			copyActions++
			if tick != 6*rate {
				t.Fatalf("copy starts at unexpected frame %d", tick)
			}
		}
		pointer := tourScreenInput("hidden", tick, total, rate)
		if pointer.PointerX < 64 || pointer.PointerX > 704 || pointer.PointerY < 68 || pointer.PointerY > 468 {
			t.Fatal("hidden-screen pointer left the visible image")
		}
		if math.Hypot(pointer.PointerX-previous.PointerX, pointer.PointerY-previous.PointerY) > 5 {
			t.Fatal("hidden-screen pointer jumped")
		}
		previous = pointer
	}
	if !reflect.DeepEqual(objects, []int{1, 3, 4, 5}) || !reflect.DeepEqual(counts, []int{4, 7, 10}) || !reflect.DeepEqual(layerKeys, []int{1, 2, 4, 3}) || copyActions != 1 {
		t.Fatalf("missing interaction: objects=%v counts=%v layers=%v copy=%d", objects, counts, layerKeys, copyActions)
	}
}

func TestTourRejectsInvalidOptionsAndStalledNavigation(t *testing.T) {
	for _, options := range []TourOptions{{}, {ScreenDuration: -time.Second, MenuDuration: time.Second}, {ScreenDuration: time.Second}} {
		if tour, err := NewTour(Config{}, options); err == nil {
			tour.Close()
			t.Fatalf("accepted invalid options %+v", options)
		}
	}
	script := newTourScript(TourOptions{ScreenDuration: time.Second, MenuDuration: time.Nanosecond}, 60)
	for i := 0; i < 3600; i++ {
		if _, _, err := script.next("menu", ""); err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err := script.next("menu", ""); err == nil {
		t.Fatal("stalled walking did not report an error")
	}
}
