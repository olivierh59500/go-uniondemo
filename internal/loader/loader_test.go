package loader

import (
	"testing"
	"time"

	"github.com/olivierh59500/go-uniondemo/assets"
	"github.com/olivierh59500/go-uniondemo/internal/screens"
)

func TestEveryScreenHasCreditsAndTimedTransition(t *testing.T) {
	for _, d := range screens.Catalog() {
		for _, rate := range []int{50, 60} {
			s, err := New(d.ID, assets.Files, rate)
			if err != nil {
				t.Fatal(err)
			}
			if len(s.lines) != 23 {
				t.Fatalf("%s: incomplete credits", d.ID)
			}
			ticks := 0
			for !s.Done() {
				s.Update()
				ticks++
			}
			elapsed := time.Duration(ticks) * time.Second / time.Duration(rate)
			if elapsed < Duration || elapsed-Duration >= time.Second/time.Duration(rate) {
				t.Fatal("loader timing is not tied to the recording duration")
			}
			s.Close()
		}
	}
}
