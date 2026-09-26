package screens

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/sprites"
	"github.com/olivierh59500/democonstructionkit/timeline"
)

func init() { factories["diskcopier"] = buildDiskCopier }

// The copy sequence is a visual performance; it never accesses a disk device.
func buildDiskCopier(s *Scene) {
	stage, textLayer := s.surface(640, 480), s.surface(640, 14)
	panel, led, lcd, raster := s.image("diskcopier.png"), s.image("led_on.png"), s.image("lcd.png"), s.image("rasters.png")
	lcdAtlas, err := sprites.NewAtlas(sprites.AtlasConfig{Image: lcd, TileW: 30, TileH: 22})
	if err != nil {
		s.err = err
		return
	}
	font := s.bitmap(s.image("font.png"), "union-diskcopier")
	var red, green [8]*ebiten.Image
	for i := range red {
		suffix := ""
		if i > 0 {
			suffix = fmt.Sprint(i)
		}
		red[i], green[i] = s.image("fontback"+suffix+".png"), s.image("fontbackg"+suffix+".png")
	}
	var copying, held bool
	var introTick, copyTick, textIndex int
	rasterBank, err := composite.NewWindowedImageBank(presets.UnionDiskCopierRasterWindows(raster))
	if err != nil {
		s.err = err
		return
	}
	s.closers = append(s.closers, rasterBank.Close)
	cues, err := timeline.NewCueRanges(presets.UnionDiskCopierCueRanges())
	if err != nil {
		s.err = err
		return
	}
	palette, err := timeline.NewSteppedEnvelope(presets.UnionDiskCopierFade())
	if err != nil || cues.Len() != len(unionCopyStages) {
		if err == nil {
			err = fmt.Errorf("disk copier cue ranges and scene stages differ")
		}
		s.err = err
		return
	}
	lcdMotion, err := motion.NewGatedWrapBank(presets.UnionDiskCopierLCDMotion())
	if err != nil {
		s.err = err
		return
	}
	x, copyX := -640.0, 0.0
	reset := func() { introTick, copyTick = 0, 0; lcdMotion.Reset(); x, copyX = -640, 0 }
	s.input = func(in Input) {
		if in.Action && !held {
			reset()
			copying = !copying
		}
		if in.Number == 10 {
			reset()
			copying = false
			textIndex = 0
		}
		held = in.Action
	}
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		textLayer.Clear()
		rasterBank.Draw(s.Canvas)
		rasterBank.Step()
		s.transform(stage, panel, 320, 200, 1, 1, 0, float64(panel.Bounds().Dx())/2, float64(panel.Bounds().Dy())/2, 1, ebiten.BlendSourceOver)
		if !copying {
			time := float64(introTick) * .5
			s.draw(textLayer, red[palette.At(time, 0, 100)], x, 0)
			if time <= 101 {
				font.Print(textLayer, unionCopyText[textIndex], 0, 0, 1, 1)
			}
			introTick++
			if time > 121 {
				textIndex = (textIndex + 1) % len(unionCopyText)
				introTick = 1
			}
			x++
			if x >= 0 {
				x = -640
			}
		} else {
			time := float64(copyTick) * .5
			start, end, line, operation := time+1, time+1, "", -1
			if index, window, active := cues.At(time); active {
				start, end = window.Start, window.End
				line, operation = unionCopyStages[index].line, unionCopyStages[index].operation
			}
			s.draw(textLayer, green[palette.At(time, start, end)], copyX, 0)
			if line != "" {
				font.Print(textLayer, line, 0, 0, 1, 1)
			}
			if operation >= 0 {
				opStart := []float64{120, 500, 760}[operation]
				leftY := 153.0
				if time > opStart+120 {
					leftY = 171
				}
				s.draw(stage, led, 280, leftY)
				s.draw(stage, led, 440, []float64{153, 189, 171}[operation])
				if (copyTick-int(opStart*2))%10 < 5 {
					s.draw(stage, led, 280, 189)
				}
				tile := lcdMotion.Frame(operation)
				s.part(stage, lcd, lcdAtlas.Region(tile), 306, 249, 1, 1)
			}
			if err := lcdMotion.StepAt(time); err != nil {
				s.err = err
				return
			}
			copyX--
			if copyX <= -640 {
				copyX = 0
			}
			if time <= 1220 {
				copyTick++
			}
		}
		s.draw(stage, textLayer, 0, 390)
		s.draw(s.Canvas, stage, 64, 28)
	}
}

var unionCopyStages = [...]struct {
	line      string
	operation int
}{
	{" PLEASE INSERT WRT-PROTECTED SOURCE-DISK ", -1},
	{"         PLEASE WAIT...  READING         ", 0},
	{"     PLEASE INSERT DESTINATION DISK      ", -1},
	{"        PLEASE WAIT... FORMATTING        ", 1},
	{"         PLEASE WAIT...  WRITING         ", 2},
	{"            ALL DONE!  ENJOY!            ", -1},
	{"   PRESS '0' TO EXIT THE COPY PROGRAM    ", -1},
}

var unionCopyText = []string{
	"        THE UNION DEMO PRESENTS:        ",
	"    COPY-PROGRAM BY 6719 AND MAD MAX     ",
	"        PRESS SPACE TO START COPY        ",
	"WHAT DO YOU THINK ABOUT THIS LITTLE COPY ",
	"    IT'S THE FIRST ONE WITH RASTERS,     ",
	"        AND MUZAK WHILE COPYING          ",
	"   AND IT TELLS YOU A LOT OF CRAP TALK   ",
	"   EXCUSE US FOR NOT USING ANY BORDER    ",
	"   OR FOR ABSENCE OF TRACKING-SPRITES    ",
	" BUT WE WANTED TO READ A WHOLE DISK-SIDE ",
	"             INTO HALF-A-MEG             ",
	" AND THIS MEANS, THAT WE NEED MEMORY !!! ",
	"DO YOU THINK THIS COPY IS A BIT TOO SLOW ",
	"REMEMBER, THIS DISC CONTAINS 900 KB DATA ",
	"          PRETTY MUCH ? YES !!           ",
	"   WE ARE USING A VERY SPECIAL FORMAT    ",
	"AND THAT TAKES IT'S TIME TO FORMAT WRITE ",
	"              AND VERIFY !!              ",
	"  AND, AFTER ALL, YOU SHALL HEAR THAT    ",
	"    FANTASTIC SOUNDTRACK COMPLETELY !    ",
	"  AND NOW SOME OF THE LATEST TEX-NEWS:   ",
	" SOME PEOPLE OF TEX ARE NOW PROFESSIONAL ",
	"             GAME-DESIGNERS              ",
	" OF COURSE WE WON'T SAY AT WHICH COMPANY ",
	"BUT WE ARE SURE YOU'LL RECOGNIZE US..... ",
	"      JUST SEE AND HEAR IF YOU'RE        ",
	"     BUYING NEW GAMES KNOWING THIS       ",
	"    YOU WON'T BE SURPRISED TO HEAR:      ",
	"  TEX WON'T CRACK ANY GAMES NO MORE!     ",
	" BUT DON'T WORRY: IF WE HAVE THE TIME,   ",
	" WE'LL CONTINUE TO MAKE DEMOS BIT BY BIT ",
	"  OK GUYS, NO MORE SPACE FOR TEXT LEFT.  ",
	"   LET'S START ANEW.  BYE, BYE FOLKS!!   ",
	"                                         ",
}
