package app

import (
	"testing"

	"github.com/olivierh59500/go-uniondemo/internal/screens"
)

func TestAllScreensLoadThroughCreditsAndReturnToMenu(t *testing.T) {
	g, err := New(Config{Muted: true, Screen: "menu"})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	for _, d := range screens.Catalog() {
		t.Run(d.ID, func(t *testing.T) {
			g.hall.PositionDoor(d.ID)
			oldStream := g.stream
			if err := g.Begin(d.ID); err != nil {
				t.Fatal(err)
			}
			if g.loader == nil || g.pending != d.ID {
				t.Fatal("missing loading transition")
			}
			if _, ok := oldStream.YMRegisters(); ok {
				t.Fatal("menu audio remained open")
			}
			for !g.loader.Done() {
				g.loader.Update()
			}
			if err := g.open(g.pending); err != nil {
				t.Fatal(err)
			}
			if g.ScreenID() != d.ID || g.scene == nil || g.loader != nil {
				t.Fatal("incorrect screen transition")
			}
			for frame := 0; frame < 5; frame++ {
				if err := g.updateAudio(); err != nil {
					t.Fatal(err)
				}
				if err := g.scene.Update(screens.Input{PointerX: 384, PointerY: 268}); err != nil {
					t.Fatal(err)
				}
			}
			screenStream := g.stream
			if err := g.open("menu"); err != nil {
				t.Fatal(err)
			}
			if g.hall.Selection() != d.ID {
				t.Fatal("return lost the selected door")
			}
			if _, ok := screenStream.YMRegisters(); ok {
				t.Fatal("screen audio remained open")
			}
		})
	}
}

func TestConfigurationAndClose(t *testing.T) {
	for _, config := range []Config{{Rate: 120}, {ScreenSeconds: -1}, {Screen: "unknown"}} {
		g, err := New(config)
		if err == nil {
			g.Close()
			t.Fatalf("accepted invalid configuration %+v", config)
		}
	}
	g, err := New(Config{Muted: true, Screen: "multiplane", Rate: 50})
	if err != nil {
		t.Fatal(err)
	}
	if g.Rate() != 50 || g.ScreenID() != "multiplane" {
		t.Fatal("initial screen/rate ignored")
	}
	if err = g.Close(); err != nil {
		t.Fatal(err)
	}
	if err = g.Close(); err != nil {
		t.Fatal(err)
	}
}
