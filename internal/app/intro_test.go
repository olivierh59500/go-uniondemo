package app

import (
	"testing"
	"time"

	"github.com/olivierh59500/democonstructionkit/sound"
	"github.com/olivierh59500/go-uniondemo/assets"
	"github.com/olivierh59500/go-uniondemo/internal/screens"
)

func TestIntroductionIsDefaultAndMenuCanBeRequested(t *testing.T) {
	for _, config := range []Config{{Muted: true}, {Muted: true, Touch: true}, {Muted: true, Screen: "intro"}} {
		g, err := New(config)
		if err != nil {
			t.Fatal(err)
		}
		if g.ScreenID() != "intro" || g.scene == nil || g.scene.Music != "audio/intro.ym" {
			t.Fatal("startup skipped the introduction or its music")
		}
		g.Close()
	}
	g, err := New(Config{Muted: true, Screen: "menu"})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	if g.ScreenID() != "menu" || g.scene != nil {
		t.Fatal("explicit menu start replayed the introduction")
	}
	if len(screens.Catalog()) != 11 {
		t.Fatal("introduction changed the hall's eleven doors")
	}
}

func TestIntroductionContinuesAtMusicEndAtEitherRate(t *testing.T) {
	data, err := assets.Files.ReadFile("audio/intro.ym")
	if err != nil {
		t.Fatal(err)
	}
	p, err := sound.Open("intro.ym", data, sound.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if p.Metadata().Duration != screens.IntroDuration {
		t.Fatal("introduction duration does not match its soundtrack")
	}
	for _, rate := range []int{50, 60} {
		g, err := New(Config{Muted: true, Rate: rate, Tour: true, ScreenSeconds: 1})
		if err != nil {
			t.Fatal(err)
		}
		end := int((screens.IntroDuration*time.Duration(rate) + time.Second - 1) / time.Second)
		g.ticks = end - 2
		if err := g.update(controls{}); err != nil {
			t.Fatal(err)
		}
		if g.ScreenID() != "intro" {
			t.Fatal("introduction ended before its music")
		}
		oldStream := g.stream
		if err := g.update(controls{}); err != nil {
			t.Fatal(err)
		}
		if g.ScreenID() != "menu" || g.scene != nil || g.loader != nil {
			t.Fatal("introduction did not lead directly to the hall")
		}
		if g.tourIndex != 0 {
			t.Fatal("introduction skipped the first screen in the tour")
		}
		if _, ok := oldStream.YMRegisters(); ok {
			t.Fatal("introduction music was not closed")
		}
		g.Close()
	}
}

func TestIntroductionCanBeSkippedAndDoesNotReplayOnReturn(t *testing.T) {
	for _, in := range []controls{{enter: true}, {back: true}, {scene: screens.Input{Action: true}}} {
		g, err := New(Config{Muted: true})
		if err != nil {
			t.Fatal(err)
		}
		if err = g.update(in); err != nil {
			t.Fatal(err)
		}
		if g.ScreenID() != "menu" {
			t.Fatal("continue input did not skip the introduction")
		}
		if err = g.open("hidden"); err != nil {
			t.Fatal(err)
		}
		if err = g.update(controls{back: true}); err != nil {
			t.Fatal(err)
		}
		if g.ScreenID() != "menu" {
			t.Fatal("return from a screen replayed the introduction")
		}
		g.Close()
	}
}
