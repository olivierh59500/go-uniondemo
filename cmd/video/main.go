// Command video records the full Union route with its own synchronized audio.
package main

import (
	"flag"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/video"
	"github.com/olivierh59500/go-uniondemo/internal/app"
)

func main() {
	config := video.Config{Output: "recordings/union-demo.mp4", Title: "The Union Demo", Width: app.Width, Height: app.Height, FPS: 60, TPS: 60, SampleRate: 48000, PosterAt: 5 * time.Second}
	config.Flags(flag.CommandLine)
	tour := app.DefaultTourOptions()
	flag.DurationVar(&tour.ScreenDuration, "screen-duration", tour.ScreenDuration, "time on each screen, excluding loading credits and walking")
	flag.DurationVar(&tour.MenuDuration, "menu-duration", tour.MenuDuration, "pause in the hall before walking to the next door")
	flag.Parse()
	if err := video.Run(config, func() (ebiten.Game, error) { return app.NewTour(app.Config{Rate: 60}, tour) }); err != nil {
		log.Fatal(err)
	}
}
