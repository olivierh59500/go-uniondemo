// Package mobile exposes the complete Union production to Android and iOS hosts.
package mobile

import (
	"github.com/hajimehoshi/ebiten/v2"
	enginemobile "github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/olivierh59500/go-uniondemo/internal/app"
)

func init() {
	game, err := app.New(app.Config{Touch: true, Rate: 60})
	if err != nil {
		panic(err)
	}
	ebiten.SetTPS(game.Rate())
	enginemobile.SetGame(game)
}

// Dummy ensures that the mobile entry point is included in the binding.
func Dummy() {}
