package screens

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/motion"
)

func init() { factories["diskcopier"] = buildDiskCopier }

// The copy sequence is a visual performance; it never accesses a disk device.
func buildDiskCopier(s *Scene) {
	stage, textLayer := s.surface(640, 480), s.surface(640, 14)
	panel, led, lcd, raster := s.image("diskcopier.png"), s.image("led_on.png"), s.image("lcd.png"), s.image("rasters.png")
	font := s.bitmap(s.image("font.png"), "union-diskcopier")
	var red, green [8]*ebiten.Image
	var strips [6]*ebiten.Image
	for i := range red {
		suffix := ""
		if i > 0 {
			suffix = fmt.Sprint(i)
		}
		red[i], green[i] = s.image("fontback"+suffix+".png"), s.image("fontbackg"+suffix+".png")
	}
	for i := range strips {
		strips[i] = s.surface(768, 32)
	}
	var copying, held bool
	var introTick, copyTick, textIndex int
	var tiles [3]float64
	rasterMotion, err := motion.NewWrapBank(motion.WrapBankConfig{
		Start: []float64{174}, Velocity: []float64{-1.5},
		Lower: &motion.WrapLimit{Boundary: -974, Restart: 174, Inclusive: true},
	})
	if err != nil {
		s.err = err
		return
	}
	x, copyX := -640.0, 0.0
	reset := func() { introTick, copyTick = 0, 0; tiles = [3]float64{}; x, copyX = -640, 0 }
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
		for i, strip := range strips {
			strip.Clear()
			s.transform(strip, raster, 0, rasterMotion.At(0)-float64(i*5), 1.3, 1, 0, 0, 0, 1, ebiten.BlendSourceOver)
			s.draw(s.Canvas, strip, 0, float64(132+i*34))
		}
		rasterMotion.Step()
		s.transform(stage, panel, 320, 200, 1, 1, 0, float64(panel.Bounds().Dx())/2, float64(panel.Bounds().Dy())/2, 1, ebiten.BlendSourceOver)
		if !copying {
			time := float64(introTick) * .5
			s.draw(textLayer, red[unionCopyFade(time, 0, 100)], x, 0)
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
			start, end, line, operation := unionCopyStage(time)
			s.draw(textLayer, green[unionCopyFade(time, start, end)], copyX, 0)
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
				tile := int(math.Floor(tiles[operation]))
				columns := max(1, lcd.Bounds().Dx()/30)
				s.part(stage, lcd, unionRegion(float64(tile%columns*30), float64(tile/columns*22), 30, 22), 306, 249, 1, 1)
			}
			for i, begin := range []float64{120, 500, 760} {
				if time > begin {
					tiles[i] += .35
					if tiles[i] > 82 {
						tiles[i] = 0
					}
				}
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

func unionCopyFade(t, start, end float64) int {
	if t < start {
		return 7
	}
	if t < start+14 {
		return max(0, min(7, 7-int((t-start)/2)))
	}
	if t >= end-16 {
		return max(0, min(7, int((t-end+16)/2)))
	}
	return 0
}

func unionCopyStage(t float64) (start, end float64, line string, operation int) {
	operation = -1
	switch {
	case t < 100:
		return 0, 100, " PLEASE INSERT WRT-PROTECTED SOURCE-DISK ", -1
	case t > 120 && t < 360:
		return 120, 360, "         PLEASE WAIT...  READING         ", 0
	case t > 380 && t < 480:
		return 380, 480, "     PLEASE INSERT DESTINATION DISK      ", -1
	case t > 500 && t < 740:
		return 500, 740, "        PLEASE WAIT... FORMATTING        ", 1
	case t > 760 && t < 1000:
		return 760, 1000, "         PLEASE WAIT...  WRITING         ", 2
	case t > 1000 && t < 1100:
		return 1000, 1100, "            ALL DONE!  ENJOY!            ", -1
	case t > 1120:
		return 1120, math.Inf(1), "   PRESS '0' TO EXIT THE COPY PROGRAM    ", -1
	default:
		return t + 1, t + 1, "", -1
	}
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
