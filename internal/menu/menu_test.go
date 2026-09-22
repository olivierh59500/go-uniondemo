package menu

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/olivierh59500/democonstructionkit/render"
	production "github.com/olivierh59500/go-uniondemo/assets"
)

func TestProductionAssetsAndDrawCadence(t *testing.T) {
	g, err := New(production.Files)
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	canvas := render.NewSurface(Width, Height)
	defer canvas.Deallocate()
	for i := 0; i < 120; i++ {
		g.Update(Input{Right: true})
	}
	before := g.state
	for i := 0; i < 3; i++ {
		g.Draw(canvas)
	}
	if g.state != before {
		t.Fatal("drawing advanced menu animation or character controls")
	}
	if err := g.Close(); err != nil {
		t.Fatal(err)
	}
	if got := g.Update(Input{Right: true, Enter: true}); got != "" || g.state != before {
		t.Fatal("closed menu accepted an input or advanced animation")
	}
	g.Draw(canvas)
}

func TestMenuRejectsMissingAndInvalidImages(t *testing.T) {
	for name, files := range map[string]fs.FS{
		"missing": fstest.MapFS{},
		"corrupt": fstest.MapFS{"menu/backdrop.png": {Data: []byte("invalid image")}},
	} {
		t.Run(name, func(t *testing.T) {
			g, err := New(files)
			if err == nil || g != nil || !strings.Contains(err.Error(), "menu/backdrop.png") {
				t.Fatalf("New = %v, %v; want named image error", g, err)
			}
		})
	}
}
