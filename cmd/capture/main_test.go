package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestCaptureAdvancesExactlyTheRequestedFrames(t *testing.T) {
	for _, frames := range []int{0, 1, 180} {
		calls := 0
		g := captureGame{options: options{frames: frames}, step: func(frame int) error {
			calls++
			if frame != calls {
				t.Fatalf("step frame = %d, expected %d", frame, calls)
			}
			return nil
		}}
		for i := 0; i < frames+4; i++ {
			if err := g.Update(); err != nil {
				t.Fatal(err)
			}
		}
		if calls != frames || g.frame != frames {
			t.Fatalf("requested %d frames, advanced %d (%d calls)", frames, g.frame, calls)
		}
		g.captured = true
		if !errors.Is(g.Update(), ebiten.Termination) {
			t.Fatal("capture does not terminate after saving its frame")
		}
	}
}

func TestCapturePropagatesUpdateFailure(t *testing.T) {
	want := errors.New("broken scene")
	g := captureGame{options: options{frames: 10}, step: func(int) error { return want }}
	if !errors.Is(g.Update(), want) {
		t.Fatal("scene failure was lost")
	}
	if g.frame != 0 {
		t.Fatal("a failed update counted as a completed frame")
	}
}

func TestSavePNGPreservesNativePixels(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 3, 2))
	source.SetRGBA(1, 1, color.RGBA{R: 31, G: 127, B: 255, A: 255})
	path := filepath.Join(t.TempDir(), "nested", "frame.png")
	if err := savePNG(path, source); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoded, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Bounds() != source.Bounds() {
		t.Fatalf("image dimensions changed: %v", decoded.Bounds())
	}
	if got := color.RGBAModel.Convert(decoded.At(1, 1)); got != source.RGBAAt(1, 1) {
		t.Fatalf("native pixel changed: %v", got)
	}
}
