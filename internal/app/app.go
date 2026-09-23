// Package app connects the Union hall, loading credits, screens and audio.
package app

import (
	"fmt"
	"io"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/sound"
	device "github.com/olivierh59500/democonstructionkit/sound/ebiten"
	"github.com/olivierh59500/democonstructionkit/sound/output"
	media "github.com/olivierh59500/go-uniondemo/assets"
	"github.com/olivierh59500/go-uniondemo/internal/loader"
	"github.com/olivierh59500/go-uniondemo/internal/menu"
	"github.com/olivierh59500/go-uniondemo/internal/screens"
)

const Width, Height = 768, 536

type Config struct {
	Screen                         string
	Muted, Touch, Tour, SkipLoader bool
	Rate                           int
	ScreenSeconds                  int
}

type Game struct {
	config             Config
	hall               *menu.Game
	scene              *screens.Scene
	loader             *loader.Screen
	pending, current   string
	stream             *sound.Stream
	player             *device.Player
	context            *output.Context
	chooser            bool
	selection          int
	ticks, tourIndex   int
	pointerX, pointerY float64
	mouseX, mouseY     int
	touchIDs           []ebiten.TouchID
	closed             bool
	mutedPCM           [7680]byte
}

func New(config Config) (*Game, error) {
	if config.Rate == 0 {
		config.Rate = 60
	}
	if config.Rate != 50 && config.Rate != 60 {
		return nil, fmt.Errorf("animation rate must be 50 or 60 Hz")
	}
	if config.ScreenSeconds == 0 {
		config.ScreenSeconds = 60
	}
	if config.ScreenSeconds < 1 {
		return nil, fmt.Errorf("screen duration must be positive")
	}
	hall, err := menu.New(media.Files)
	if err != nil {
		return nil, err
	}
	g := &Game{config: config, hall: hall, pointerX: Width / 2, pointerY: Height / 2}
	id := config.Screen
	if id == "" {
		id = "intro"
	}
	if err = g.open(id); err != nil {
		g.Close()
		return nil, err
	}
	return g, nil
}

func (g *Game) open(id string) error {
	var next *screens.Scene
	var err error
	music := "audio/menu.ym"
	var loopStartFrame int64
	if id != "menu" {
		next, err = screens.New(id, media.Files)
		if err != nil {
			return err
		}
		music = next.Music
		loopStartFrame = next.Descriptor.MusicLoopStartFrame
	}
	if g.scene != nil {
		g.scene.Close()
	}
	g.scene = next
	if g.loader != nil {
		g.loader.Close()
		g.loader = nil
	}
	g.current = id
	g.pending = ""
	g.ticks = 0
	return g.setMusic(music, id != "intro", loopStartFrame)
}

// Begin preserves the hall position and presents the selected screen's credits.
func (g *Game) Begin(id string) error {
	if id == "intro" {
		return g.open(id)
	}
	if _, ok := screens.Find(id); !ok {
		return fmt.Errorf("unknown screen %q", id)
	}
	if g.config.SkipLoader {
		return g.open(id)
	}
	next, err := loader.New(id, media.Files, g.config.Rate)
	if err != nil {
		return err
	}
	if g.loader != nil {
		g.loader.Close()
	}
	if g.scene != nil {
		g.scene.Close()
		g.scene = nil
	}
	g.loader = next
	g.pending = id
	g.ticks = 0
	return g.setMusic("audio/loader.ogg", false, 0)
}

func (g *Game) setMusic(name string, loop bool, loopStartFrame int64) error {
	if g.player != nil {
		g.player.Close()
		g.player = nil
		g.stream = nil
	}
	if g.stream != nil {
		g.stream.Close()
		g.stream = nil
	}
	if name == "" {
		return nil
	}
	data, err := media.Files.ReadFile(name)
	if err != nil {
		return err
	}
	g.stream, err = sound.Open(name, data, sound.Options{SampleRate: 48000, Loop: loop, LoopStartFrame: loopStartFrame})
	return err
}

func (g *Game) updateAudio() error {
	if g.stream == nil {
		return nil
	}
	if g.config.Muted {
		if g.player != nil {
			g.player.SetVolume(0)
		}
		if g.player == nil {
			_, err := io.ReadFull(g.stream, g.mutedPCM[:48000/g.config.Rate*8])
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}
		}
	} else if g.player == nil {
		if g.context == nil {
			g.context = output.CurrentContext()
			if g.context == nil {
				g.context = output.NewContext(48000)
			}
		}
		p, err := device.NewOutputPlayer(g.context, g.stream)
		if err != nil {
			return err
		}
		g.player = p
		g.player.SetVolume(.7)
		g.player.Play()
	} else {
		g.player.SetVolume(.7)
	}
	if g.scene != nil {
		if registers, ok := g.stream.YMRegisters(); ok {
			for i := 0; i < 3; i++ {
				g.scene.VoiceVolumes[i] = registers[8+i] & 31
			}
		}
	}
	return nil
}

func (g *Game) Update() error {
	if g.closed {
		return ebiten.Termination
	}
	return g.update(g.controls())
}

func (g *Game) introComplete() bool {
	if g.player != nil {
		// Wait for audible playback, including the device's buffered samples.
		return g.player.Position() >= screens.IntroDuration
	}
	return g.Position() >= screens.IntroDuration
}

// update processes one fixed-rate input snapshot independently of its source.
func (g *Game) update(in controls) error {
	if g.closed {
		return ebiten.Termination
	}
	if in.fullscreen {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if in.mute {
		g.config.Muted = !g.config.Muted
	}
	if in.choose {
		g.chooser = !g.chooser
	}
	if in.back {
		g.chooser = false
		return g.open("menu")
	}
	if g.chooser {
		catalog := screens.Catalog()
		if in.previous {
			g.selection = (g.selection + len(catalog) - 1) % len(catalog)
		}
		if in.next {
			g.selection = (g.selection + 1) % len(catalog)
		}
		if in.enter {
			g.chooser = false
			g.hall.PositionDoor(catalog[g.selection].ID)
			return g.Begin(catalog[g.selection].ID)
		}
		return g.updateAudio()
	}
	if err := g.updateAudio(); err != nil {
		return err
	}
	g.ticks++
	if g.loader != nil {
		g.loader.Update()
		if g.loader.Done() || in.enter {
			return g.open(g.pending)
		}
		return nil
	}
	if g.current == "intro" {
		if in.enter || in.scene.Action || g.introComplete() {
			return g.open("menu")
		}
		return g.scene.Update(in.scene)
	}
	if g.scene == nil {
		id := g.hall.Update(menu.Input{Left: in.scene.Left, Right: in.scene.Right, Up: in.scene.Up, Down: in.scene.Down, Enter: in.enter})
		if id != "" {
			return g.Begin(id)
		}
		if g.config.Tour && g.ticks >= g.config.Rate*4 {
			catalog := screens.Catalog()
			id = catalog[g.tourIndex%len(catalog)].ID
			g.hall.PositionDoor(id)
			return g.Begin(id)
		}
	} else {
		if err := g.scene.Update(in.scene); err != nil {
			return err
		}
		if g.config.Tour && g.ticks >= g.config.Rate*g.config.ScreenSeconds {
			g.tourIndex++
			return g.open("menu")
		}
	}
	return nil
}

func (g *Game) Draw(dst *ebiten.Image) {
	if g.loader != nil {
		g.loader.Draw(dst)
	} else if g.scene != nil {
		g.scene.Draw(dst)
	} else {
		g.hall.Draw(dst)
	}
	if g.config.Touch {
		g.drawPad(dst)
	}
	if g.chooser {
		g.drawChooser(dst)
	}
}
func (g *Game) Layout(int, int) (int, int) { return Width, Height }
func (g *Game) ScreenID() string {
	if g.loader != nil {
		return "loading:" + g.pending
	}
	return g.current
}
func (g *Game) Rate() int { return g.config.Rate }
func (g *Game) Position() time.Duration {
	return time.Duration(g.ticks) * time.Second / time.Duration(g.config.Rate)
}
func (g *Game) Close() error {
	if g.closed {
		return nil
	}
	g.closed = true
	if g.player != nil {
		g.player.Close()
		g.player = nil
		g.stream = nil
	}
	if g.stream != nil {
		g.stream.Close()
		g.stream = nil
	}
	if g.scene != nil {
		g.scene.Close()
		g.scene = nil
	}
	if g.loader != nil {
		g.loader.Close()
		g.loader = nil
	}
	return g.hall.Close()
}
