// Command capture exports a deterministic native-resolution frame without opening
// an audio device. Music-driven animation is advanced by the YM synthesizer.
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/sound"
	"github.com/olivierh59500/go-uniondemo/assets"
	"github.com/olivierh59500/go-uniondemo/internal/loader"
	"github.com/olivierh59500/go-uniondemo/internal/menu"
	"github.com/olivierh59500/go-uniondemo/internal/screens"
)

type options struct {
	screen, output, door  string
	frames, rate, number  int
	pointerMotion, action bool
}

func (o options) validate() error {
	if o.frames < 0 {
		return fmt.Errorf("frames must be nonnegative")
	}
	if o.rate != 50 && o.rate != 60 {
		return fmt.Errorf("rate must be 50 or 60")
	}
	if o.output == "" {
		return fmt.Errorf("output path must not be empty")
	}
	if o.number < 0 || o.number > 10 {
		return fmt.Errorf("number must be between 0 and 10")
	}
	if o.door != "" {
		if _, ok := screens.Find(o.door); !ok {
			return fmt.Errorf("unknown door %q", o.door)
		}
	}
	if o.screen != "menu" {
		if _, ok := screens.Find(strings.TrimPrefix(o.screen, "loader:")); !ok {
			return fmt.Errorf("unknown screen %q", o.screen)
		}
	}
	return nil
}

type captureGame struct {
	options              options
	width, height, frame int
	step                 func(int) error
	draw                 func(*ebiten.Image)
	close                func() error
	captured             bool
	err                  error
}

func newCapture(o options) (*captureGame, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}
	g := &captureGame{options: o, width: menu.Width, height: menu.Height}
	switch {
	case o.screen == "menu":
		scene, err := menu.New(assets.Files)
		if err != nil {
			return nil, err
		}
		scene.PositionDoor(o.door)
		g.step = func(int) error { scene.Update(menu.Input{}); return nil }
		g.draw, g.close = scene.Draw, scene.Close
	case strings.HasPrefix(o.screen, "loader:"):
		scene, err := loader.New(strings.TrimPrefix(o.screen, "loader:"), assets.Files, o.rate)
		if err != nil {
			return nil, err
		}
		g.step = func(int) error { scene.Update(); return nil }
		g.draw, g.close = scene.Draw, scene.Close
	default:
		scene, err := screens.New(o.screen, assets.Files)
		if err != nil {
			return nil, err
		}
		g.width, g.height = scene.Descriptor.Width, scene.Descriptor.Height
		g.draw, g.close = scene.Draw, scene.Close
		var music *sound.Stream
		var samples []byte
		if o.screen == "delta" {
			data, readErr := assets.Files.ReadFile(scene.Music)
			if readErr != nil {
				scene.Close()
				return nil, readErr
			}
			music, err = sound.Open(scene.Music, data, sound.Options{SampleRate: 48000, Loop: true, BlockFrames: 48000 / o.rate})
			if err != nil {
				scene.Close()
				return nil, err
			}
			samples = make([]byte, 48000/o.rate*8)
			g.close = func() error { return errors.Join(music.Close(), scene.Close()) }
		}
		g.step = func(frame int) error {
			if music != nil {
				if _, err := io.ReadFull(music, samples); err != nil {
					return err
				}
				if registers, ok := music.YMRegisters(); ok {
					for voice := range scene.VoiceVolumes {
						scene.VoiceVolumes[voice] = registers[8+voice]
					}
				}
			}
			input := screens.Input{PointerX: float64(g.width) / 2, PointerY: float64(g.height) / 2}
			if o.pointerMotion {
				t := float64(frame) / float64(o.rate)
				input.PointerX += 220 * math.Sin(t*2)
				input.PointerY += 140 * math.Cos(t*3)
			}
			if frame == 1 {
				input.Number, input.Action = o.number, o.action
			}
			return scene.Update(input)
		}
	}
	return g, nil
}

func (g *captureGame) Update() error {
	if g.captured {
		return ebiten.Termination
	}
	if g.frame < g.options.frames {
		if err := g.step(g.frame + 1); err != nil {
			return err
		}
		g.frame++
	}
	return nil
}

func (g *captureGame) Draw(dst *ebiten.Image) {
	g.draw(dst)
	if !g.captured && g.frame == g.options.frames {
		pixels := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
		dst.ReadPixels(pixels.Pix)
		g.err = savePNG(g.options.output, pixels)
		g.captured = true
	}
}

func (g *captureGame) Layout(_, _ int) (int, int) { return g.width, g.height }

func savePNG(path string, pixels image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".union-frame-*.png")
	if err != nil {
		return err
	}
	temporary := file.Name()
	defer os.Remove(temporary)
	if err = png.Encode(file, pixels); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}

func run(o options) error {
	g, err := newCapture(o)
	if err != nil {
		return err
	}
	defer g.close()
	ebiten.SetWindowSize(g.width, g.height)
	ebiten.SetWindowTitle("The Union Demo — Frame Capture")
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetVsyncEnabled(false)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetScreenClearedEveryFrame(true)
	if err := ebiten.RunGame(g); err != nil && !errors.Is(err, ebiten.Termination) {
		return err
	}
	if !g.captured {
		return fmt.Errorf("capture window closed before frame %d", o.frames)
	}
	return g.err
}

func main() {
	var o options
	flag.StringVar(&o.screen, "screen", "menu", "menu, loader:<screen>, or a screen identifier")
	flag.IntVar(&o.frames, "frames", 180, "exact number of updates before capture")
	flag.StringVar(&o.output, "output", "capture.png", "output PNG path")
	flag.IntVar(&o.rate, "rate", 60, "animation and music ticks per second (50 or 60)")
	flag.StringVar(&o.door, "door", "", "place the menu character at this door")
	flag.BoolVar(&o.pointerMotion, "pointer-motion", false, "move the pointer along a deterministic curve")
	flag.IntVar(&o.number, "number", 0, "press this number on the first update; 10 selects the zero key")
	flag.BoolVar(&o.action, "action", false, "press the action button on the first update")
	flag.Parse()
	if err := run(o); err != nil {
		log.Fatal(err)
	}
}
