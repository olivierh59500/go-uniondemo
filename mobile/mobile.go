// Package mobile exposes the complete Union production to Android and iOS hosts.
package mobile

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	enginemobile "github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/olivierh59500/go-uniondemo/internal/app"
)

type host struct {
	game         *app.Game
	tourRequests chan int
	touring      bool
	lastScreen   string
}

func (h *host) Update() error {
	select {
	case seconds := <-h.tourRequests:
		if err := h.game.SetTour(seconds); err != nil {
			return err
		}
		h.touring = true
	default:
	}
	if err := h.game.Update(); err != nil {
		return err
	}
	if h.touring {
		if screen := h.game.ScreenID(); screen != h.lastScreen {
			log.Printf("union_tour screen=%s", screen)
			h.lastScreen = screen
		}
	}
	return nil
}

func (h *host) Draw(dst *ebiten.Image) { h.game.Draw(dst) }
func (h *host) Layout(width, height int) (int, int) {
	return h.game.Layout(width, height)
}

var gameHost *host

func init() {
	game, err := app.New(app.Config{Touch: true, Rate: 60})
	if err != nil {
		panic(err)
	}
	gameHost = &host{game: game, tourRequests: make(chan int, 1)}
	ebiten.SetTPS(game.Rate())
	enginemobile.SetGame(gameHost)
}

// ConfigureTour queues a full-screen tour before Android starts drawing.
// Zero is not a tour; positive values from one to 600 seconds are accepted.
func ConfigureTour(screenSeconds int) bool {
	if gameHost == nil || screenSeconds < 1 || screenSeconds > 600 {
		return false
	}
	select {
	case <-gameHost.tourRequests:
	default:
	}
	gameHost.tourRequests <- screenSeconds
	return true
}

// Dummy ensures that the mobile entry point is included in the binding.
func Dummy() {}
