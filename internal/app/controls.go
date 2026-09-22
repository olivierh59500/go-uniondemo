package app

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/olivierh59500/go-uniondemo/internal/screens"
)

type controls struct {
	scene                                                 screens.Input
	enter, back, choose, mute, fullscreen, previous, next bool
}
type button struct {
	label string
	rect  image.Rectangle
}

var buttons = []button{
	{"LEFT", image.Rect(12, 482, 78, 526)}, {"UP", image.Rect(82, 462, 140, 495)},
	{"DOWN", image.Rect(82, 498, 140, 531)}, {"RIGHT", image.Rect(144, 482, 210, 526)},
	{"ENTER", image.Rect(342, 482, 430, 526)}, {"MENU", image.Rect(436, 482, 524, 526)},
	{"SCREENS", image.Rect(530, 482, 634, 526)}, {"ACTION", image.Rect(640, 482, 752, 526)},
}

func (g *Game) controls() controls {
	just := inpututil.IsKeyJustPressed
	x, y := ebiten.CursorPosition()
	if (x != g.mouseX || y != g.mouseY) && image.Pt(x, y).In(image.Rect(0, 0, Width, Height)) {
		g.pointerX, g.pointerY = float64(x), float64(y)
	}
	g.mouseX, g.mouseY = x, y
	in := controls{scene: screens.Input{Left: ebiten.IsKeyPressed(ebiten.KeyArrowLeft), Right: ebiten.IsKeyPressed(ebiten.KeyArrowRight), Up: ebiten.IsKeyPressed(ebiten.KeyArrowUp), Down: ebiten.IsKeyPressed(ebiten.KeyArrowDown), Action: just(ebiten.KeySpace), PointerX: g.pointerX, PointerY: g.pointerY, PointerDown: ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)}, enter: just(ebiten.KeyEnter) || just(ebiten.KeySpace), back: just(ebiten.KeyEscape) || just(ebiten.KeyBackspace), choose: just(ebiten.KeyTab), mute: just(ebiten.KeyM), fullscreen: just(ebiten.KeyF), previous: just(ebiten.KeyArrowUp), next: just(ebiten.KeyArrowDown)}
	for i, key := range []ebiten.Key{ebiten.KeyDigit1, ebiten.KeyDigit2, ebiten.KeyDigit3, ebiten.KeyDigit4, ebiten.KeyDigit5, ebiten.KeyDigit6, ebiten.KeyDigit7, ebiten.KeyDigit8, ebiten.KeyDigit9, ebiten.KeyDigit0} {
		if just(key) {
			in.scene.Number = i + 1
		}
	}
	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])
	apply := func(x, y int, pressed bool) {
		g.pointerX, g.pointerY = float64(x), float64(y)
		in.scene.PointerX, in.scene.PointerY = g.pointerX, g.pointerY
		in.scene.PointerDown = true
		if !g.config.Touch {
			return
		}
		for _, b := range buttons {
			if !image.Pt(x, y).In(b.rect) {
				continue
			}
			switch b.label {
			case "LEFT":
				in.scene.Left = true
			case "RIGHT":
				in.scene.Right = true
			case "UP":
				in.scene.Up = true
				in.previous = in.previous || pressed
			case "DOWN":
				in.scene.Down = true
				in.next = in.next || pressed
			case "ENTER":
				in.enter = in.enter || pressed
			case "MENU":
				in.back = in.back || pressed
			case "SCREENS":
				in.choose = in.choose || pressed
			case "ACTION":
				in.scene.Action = in.scene.Action || pressed
			}
		}
	}
	for _, id := range g.touchIDs {
		x, y := ebiten.TouchPosition(id)
		apply(x, y, inpututil.TouchPressDuration(id) == 1)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		apply(x, y, inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft))
	}
	return in
}
func (g *Game) drawPad(dst *ebiten.Image) {
	for _, b := range buttons {
		vector.DrawFilledRect(dst, float32(b.rect.Min.X), float32(b.rect.Min.Y), float32(b.rect.Dx()), float32(b.rect.Dy()), color.RGBA{20, 24, 40, 215}, false)
		vector.StrokeRect(dst, float32(b.rect.Min.X), float32(b.rect.Min.Y), float32(b.rect.Dx()), float32(b.rect.Dy()), 1, color.RGBA{100, 120, 180, 255}, false)
		ebitenutil.DebugPrintAt(dst, b.label, b.rect.Min.X+(b.rect.Dx()-len(b.label)*6)/2, b.rect.Min.Y+(b.rect.Dy()-16)/2)
	}
}
func (g *Game) drawChooser(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, 132, 54, 504, 410, color.RGBA{8, 12, 24, 245}, false)
	ebitenutil.DebugPrintAt(dst, "THE UNION DEMO - SCREENS", 176, 76)
	for i, d := range screens.Catalog() {
		marker := "  "
		if i == g.selection {
			marker = "> "
		}
		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%s%02d  %s", marker, i+1, strings.ReplaceAll(d.Title, "—", "-")), 160, 108+i*26)
	}
	ebitenutil.DebugPrintAt(dst, "UP/DOWN: SELECT    ENTER: OPEN    ESC: MENU", 154, 426)
}
