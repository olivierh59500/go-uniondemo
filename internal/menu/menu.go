package menu

import (
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/olivierh59500/democonstructionkit/assets"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

const (
	Width  = 768
	Height = 536
)

const welcomeText = "     THE UNION DEMO     THE CAREBEARS - TEX - TNT CREW - THE REPLICANTS - DELTA FORCE - LEVEL 16     USE THE ARROW KEYS OR THE DIRECTION PAD TO WALK CHARLY TO A DOOR. PRESS ENTER OR THE ENTER BUTTON TO ENJOY THE SCREEN. PRESS ESCAPE OR THE MENU BUTTON TO RETURN TO THE UNION HALL.     "

// Game owns the menu's surfaces and shared DCK effects. The caller supplies the
// fixed 60 Hz update clock, audio playback, and transitions into selected screens.
type Game struct {
	state
	store                                                       *assets.Store
	backdrop, bannerLogo, bannerPattern, bar, panorama, charley *ebiten.Image
	bannerSurface, hallSurface, scrollSurface                   *ebiten.Image
	text                                                        *scrolling.Scrolling
	closed                                                      bool
}

// New loads all menu resources before returning a playable menu.
func New(files fs.FS) (_ *Game, err error) {
	g := &Game{state: newState(), store: assets.New(files)}
	defer func() {
		if err != nil {
			_ = g.Close()
		}
	}()
	resources := []struct {
		name string
		dst  **ebiten.Image
		min  image.Point
	}{
		{"backdrop.png", &g.backdrop, image.Pt(640, 246)},
		{"bannerlogo.png", &g.bannerLogo, image.Pt(1, 1)},
		{"banner4.png", &g.bannerPattern, image.Pt(640, 58)},
		{"bar.png", &g.bar, image.Pt(1, 1)},
		{"panorama.png", &g.panorama, image.Pt(768, 32)},
		{"charley.png", &g.charley, image.Pt(80, 102)},
	}
	for _, r := range resources {
		name := "menu/" + r.name
		img, loadErr := g.store.Image(name)
		if loadErr != nil {
			return nil, fmt.Errorf("menu: %w", loadErr)
		}
		if img.Bounds().Dx() < r.min.X || img.Bounds().Dy() < r.min.Y {
			return nil, fmt.Errorf("menu: %s is smaller than %dx%d", name, r.min.X, r.min.Y)
		}
		*r.dst, loadErr = g.store.Texture(name)
		if loadErr != nil {
			return nil, fmt.Errorf("menu: %w", loadErr)
		}
	}
	if (g.charley.Bounds().Dx()/80)*(g.charley.Bounds().Dy()/102) < 8 {
		return nil, fmt.Errorf("menu: character atlas needs eight 80x102 frames")
	}
	atlas, err := g.store.Texture("menu/fontsTexOut2.png")
	if err != nil {
		return nil, fmt.Errorf("menu: %w", err)
	}
	grid, err := presets.BitmapFont("union-menu", atlas, ebiten.FilterNearest)
	if err != nil {
		return nil, fmt.Errorf("menu: font: %w", err)
	}
	g.text, err = grid.Scrolling(welcomeText)

	if err != nil {
		return nil, fmt.Errorf("menu: scrolling: %w", err)
	}
	g.bannerSurface = render.NewSurface(640, 58)
	g.hallSurface = render.NewSurface(640, 246)
	g.scrollSurface = render.NewSurface(768, 34)
	return g, nil
}

// Update advances one 60 Hz tick and returns a door only for an enter action.
func (g *Game) Update(in Input) string {
	if g.closed {
		return ""
	}
	return g.state.update(in)
}

// Selection returns the screen at Charly's current door, or an empty string.
func (g *Game) Selection() string { return g.state.selection() }

// PositionDoor places Charly at a door without resetting menu animation. Unknown
// screen identifiers leave the current position unchanged.
func (g *Game) PositionDoor(id string) { g.state.positionDoor(id) }

// Draw composes the current state without advancing any animation or controls.
func (g *Game) Draw(dst *ebiten.Image) {
	if g.closed || dst == nil {
		return
	}
	dst.Fill(color.Black)
	bar := ebiten.DrawImageOptions{}
	bar.GeoM.Scale(77, 1)
	composite.Instance{Image: g.bar, Options: bar}.Draw(dst)

	g.bannerSurface.Clear()
	drawImage(g.bannerSurface, g.bannerPattern, float64(g.bannerX), 0)
	logo := ebiten.DrawImageOptions{}
	logo.GeoM.Translate(-float64(g.bannerLogo.Bounds().Dx())/2, -float64(g.bannerLogo.Bounds().Dy())/2)
	logo.GeoM.Scale(1, g.logoScale)
	logo.GeoM.Translate(320, 29)
	composite.Instance{Image: g.bannerLogo, Options: logo}.Draw(g.bannerSurface)
	drawImage(dst, g.bannerSurface, 64, 58)

	g.hallSurface.Fill(hallColor(g.colorIndex))
	drawImage(g.hallSurface, g.backdrop, float64(g.backX), 0)
	columns := g.charley.Bounds().Dx() / 80
	frame := image.Rect((g.charleyFrame%columns)*80, (g.charleyFrame/columns)*102,
		(g.charleyFrame%columns+1)*80, (g.charleyFrame/columns+1)*102)
	character := ebiten.DrawImageOptions{}
	character.GeoM.Scale(float64(g.facing), 1)
	x := 300
	if g.facing < 0 {
		x += 80
	}
	character.GeoM.Translate(float64(x), float64(g.charleyY))
	composite.Instance{Image: g.charley, Source: &frame, Options: character}.Draw(g.hallSurface)
	drawImage(dst, g.hallSurface, 64, 118)

	g.scrollSurface.Clear()
	drawImage(g.scrollSurface, g.panorama, float64(g.panoramaX), 0)
	if g.coverWidth > 0 {
		vector.DrawFilledRect(g.scrollSurface, 0, 0, float32(g.coverWidth), 34, color.RGBA{G: 32, B: 32, A: 255}, false)
	}
	distance := float64(g.tick) * 5
	textState := scrolling.IdentityState()
	textState.X = 768 - distance
	textState.First = max(0, int(math.Floor((distance-768)/64)))
	textState.End = textState.First + 14
	textState.Cycle = true
	textState.Time = float64(g.tick) / 60
	textState.Position = distance
	g.text.DrawAt(g.scrollSurface, textState)
	drawImage(dst, g.scrollSurface, 0, 396)
}

func drawImage(dst, src *ebiten.Image, x, y float64) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	composite.Instance{Image: src, Options: op}.Draw(dst)
}

// Close releases the menu's textures. It is safe to call more than once.
func (g *Game) Close() error {
	if g.closed {
		return nil
	}
	g.closed = true
	g.text = nil
	for _, surface := range []*ebiten.Image{g.bannerSurface, g.hallSurface, g.scrollSurface} {
		if surface != nil {
			surface.Deallocate()
		}
	}
	return g.store.Close()
}
