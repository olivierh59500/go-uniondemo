# The Union Demo

A native Go/Ebitengine production built with Democonstructionkit (DCK).
The introduction opens with a scrolling background, waving text and YM music,
then leads to the Union hall and its eleven screens with animated loading credits.
Graphics, bitmap fonts and music are embedded; playback works offline.

## Run

Go 1.26 or newer is required. Dependencies are pinned in `go.mod`, including
YM Player v1.0.0 and a published DCK version.

```sh
go run ./cmd/uniondemo
go run ./cmd/uniondemo -screen menu
go run ./cmd/uniondemo -list
go run ./cmd/uniondemo -screen multiplane
go run ./cmd/uniondemo -tour -seconds 60
go run ./cmd/uniondemo -touch
```

Startup includes the introduction on desktop and Android. Enter, Space or the
ENTER/MENU touch button continues immediately; otherwise the menu opens after
the complete 38.42-second introduction tune. `-screen menu` starts directly in
the hall. Returning from a screen does not replay the introduction.

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

The introduction uses *Think Twice* in YM format. Ten screen soundtracks and
the menu also use YM playback through DCK. Beat Dis and Wow use sampled YM tracks. The hidden screen plays
Mad Max's *Thalion Forever (Feed Me Max)* as losslessly compressed PCM, retaining
its introduction and continuous musical loop. The loading transition plays its
short recorded cue once. All music is opened through
`sound.Open(filename, data, options)`, including gzip-wrapped recordings. DCK
selects the decoder and maintains the requested loop start; the demo and capture
tool contain no decoder selection or PCM conversion. Playback remains entirely
in Go.

Reusable composition, scrolling, raster deformation, sprite and mesh effects
come from DCK. Each font supplies its own atlas dimensions and ordering. The
Delta Force balls read synchronized YM voice registers from the active stream.
The credits loader now uses DCK's `timeline.CueClock` for its exact recorded
duration and reveal clock. Seven captures around its end boundary match the
previous loader pixel for pixel.
Disk Copier now selects its copy stages through `timeline.CueRanges` and its
eight palette images through `timeline.SteppedEnvelope`. The messages, controls
and LED/LCD images remain scene data. Nine captures through its early and late
copy operations match the previous screen pixel for pixel.
Scene clocks advance only in Update; rendering does not advance animation.
Wow and Replicants use `composite.RasterOverlay` for their masked raster
material. Their opposite phase directions and inclusive wrap boundaries stay
editable in DCK. Replicants arranges its letter sprites in a resting row and
uses DCK's `motion.CuedFormation` with `sprites.Group` to sequence staggered
horizontal, vertical, arcing and compressing trajectories. The six old
Replicants pixel captures predate this movement correction; its raster
handoffs remain unchanged. The two bouncing raster banks now use DCK's
`sprites.Train` with `motion.BounceBank`; ten frames around the bounce
boundaries match the previous implementation pixel for pixel.
Beat Dis, Wow, TNT2 and Level 16 now share `motion.WrapBank` for their
independent scrolling layers. Thirty-five checkpoints, including TNT2 control
changes and strict/inclusive wrap boundaries, match the previous screens pixel
for pixel.
Beat Dis also uses `motion.HarmonicFormation` and `sprites.Group` for its eight
letters. A phase table preserves their nonuniform spacing while separate clocks
control the common horizontal wobble and the individual orbits.

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

## Complete video recording

With FFmpeg installed and a native display available:

```sh
go run ./cmd/video
go run ./cmd/video -output recordings/union-preview.mp4 -duration 10s
```

The default export visits all eleven screens for one minute each, including the
complete introduction, loading credits, Charly's walk between doors, and four
seconds in the hall on every return. The route ends after the final return to
the hall, in about fourteen minutes. It demonstrates all five TNT3 objects,
several star counts and layer controls, the hidden screen's pointer trails,
and the animated disk-copy sequence.

The recording contains only the 768 × 536 canvas at 60 frames per second and
the production's own stereo audio at 48 kHz. Animation and music share the same
offline clock, independently of export speed. The command creates an MP4,
a PNG poster from the introduction, and a JSON report with chapter timings in
the locally excluded `recordings/` directory. `-screen-duration`,
`-menu-duration` and `-poster-at` adjust the presentation; `-duration` limits a
preview without changing the complete route.

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
