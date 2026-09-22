# The Union Demo

A native Go/Ebitengine production built with Democonstructionkit (DCK).
The Union hall connects eleven screens through animated loading credits.
Graphics, bitmap fonts and music are embedded; playback works offline.

## Run

Go 1.26 or newer is required. Dependencies are pinned in `go.mod`, including
YM Player v1.0.0 and a published DCK version.

```sh
go run ./cmd/uniondemo
go run ./cmd/uniondemo -list
go run ./cmd/uniondemo -screen multiplane
go run ./cmd/uniondemo -tour -seconds 60
go run ./cmd/uniondemo -touch
```

The fixed animation rate defaults to 60 Hz independently of display refresh.
Use `-rate 50` for 50 Hz playback, `-mute` for silent playback, or
`-skip-loader` to open screens immediately.

## Screens

| ID | Screen |
| --- | --- |
| `beatdis` | The Carebears — Beat Dis |
| `delta` | Delta Force — Spher-I-Cool |
| `tnt3` | TNT Crew — Solid 3D |
| `wow` | The Carebears — Wow Scroller |
| `hidden` | Hidden Screen |
| `starballs` | TNT Crew — Starballs |
| `replicants` | The Replicants |
| `tnt2` | TNT Crew — Layered Scroller |
| `level16` | Level 16 |
| `multiplane` | The Carebears — Multi-Plane 3D Scroller |
| `diskcopier` | TEX — Disk Copier |

## Controls

- Arrow keys: move Charly in the hall, or adjust the current effect.
- Enter/Space: enter a door; Enter also skips the loading credits.
- Escape/Backspace: return to the same place in the hall.
- Tab: screen chooser; Up/Down and Enter select an entry.
- F: fullscreen. M: mute.
- TNT3: 1–5 select the five solid objects; Space cycles them.
- Starballs: 1–9 and 0 adjust the star count; Up/Down or Action also work.
- TNT2: 1–4 select motion; arrows control direction/speed.
- Replicants: Left/Right adjust scrolling speed.
- Disk Copier: Space/Action starts the animated copy sequence; 0 stops it.

The touch overlay provides directions, ENTER, MENU, SCREENS and ACTION.
The disk-copy screen is a visual sequence and does not access disks or user files.

## Audio and effects

All eleven screen soundtracks and the menu use YM playback through DCK and
YM Player v1.0.0. The loading transition uses its short recorded cue. The hidden
screen uses the Union track *Think Twice*. Beat Dis and Wow use sampled YM tracks.

Reusable composition, scrolling, raster deformation, sprite and mesh effects
come from DCK. Each font supplies its own atlas dimensions and ordering. The
Delta Force balls read synchronized YM voice registers from the active stream.
Scene clocks advance only in Update; rendering does not advance animation.

## Validation and frame captures

```sh
go test -race ./...
go vet ./...
go run ./cmd/capture -screen menu -frames 180 -output /tmp/union-menu.png
go run ./cmd/capture -screen delta -frames 600 -output /tmp/union-delta.png
go run ./cmd/capture -screen tnt3 -number 3 -frames 240 -output /tmp/union-sphere.png
```

Graphics checks require a native display. The capture command renders at native
resolution without opening an audio device, advancing music-driven animation
with the YM synthesizer. Its frame count is the number of updates after the
initial screen image. `-pointer-motion` exercises the hidden screen's trails.

## Android

The `mobile` package uses the complete production with touch controls. The
Android wrapper targets ARM64, Android 6 or newer, and uses SDK 36 / JDK 17.

```sh
./scripts/run-android.sh --build-only
./scripts/run-android.sh
```

The second command installs the debug build on the authorized connected device.
Set `ANDROID_SERIAL` when several devices are connected. The application ID is
`com.olivierh.uniondemo`.

## Credits

The Union's original artists, programmers and musicians retain credit for their
production. The loading cards preserve the individual screen credits. This
repository contains the native application and its assets.
