package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/go-uniondemo/internal/app"
	"github.com/olivierh59500/go-uniondemo/internal/screens"
)

func main() {
	var config app.Config
	var list bool
	flag.StringVar(&config.Screen, "screen", "menu", "Start at menu or a named screen")
	flag.BoolVar(&config.Muted, "mute", false, "Disable sound output")
	flag.BoolVar(&config.Touch, "touch", false, "Show the touch controls")
	flag.BoolVar(&config.Tour, "tour", false, "Visit every screen with returns to the menu")
	flag.BoolVar(&config.SkipLoader, "skip-loader", false, "Open selected screens immediately")
	flag.IntVar(&config.Rate, "rate", 60, "Animation ticks per second: 50 or 60")
	flag.IntVar(&config.ScreenSeconds, "seconds", 60, "Time on each screen during the tour")
	flag.BoolVar(&list, "list", false, "List the available screens")
	flag.Parse()
	if list {
		for _, d := range screens.Catalog() {
			fmt.Printf("%-14s %s\n", d.ID, d.Title)
		}
		return
	}
	game, err := app.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer game.Close()
	ebiten.SetTPS(game.Rate())
	ebiten.SetWindowTitle("The Union Demo")
	ebiten.SetWindowSize(1152, 804)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err = ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
