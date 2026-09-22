package app

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/olivierh59500/go-uniondemo/assets"
)

// Recorded screen music must repeat, while the loading cue must play once.
func TestRecordedMusicLoopPolicy(t *testing.T) {
	g, err := New(Config{Muted: true})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	if err := g.setMusic("audio/loader.ogg", false, 0); err != nil {
		t.Fatal(err)
	}
	cycle, err := io.ReadAll(g.stream)
	if err != nil {
		t.Fatal(err)
	}
	if len(cycle) < 48000*8 {
		t.Fatal("loading cue ended before its first second")
	}
	if err := g.setMusic("audio/loader.ogg", true, 0); err != nil {
		t.Fatal(err)
	}
	got := make([]byte, len(cycle)+4096)
	if _, err := io.ReadFull(g.stream, got); err != nil {
		t.Fatal("recorded screen music did not loop:", err)
	}
	if !bytes.Equal(got[:len(cycle)], cycle) || !bytes.Equal(got[len(cycle):], cycle[:4096]) {
		t.Fatal("recorded loop changed the PCM or inserted silence")
	}
}

func TestHiddenMusicPreservesIntroductionAndLoopSeam(t *testing.T) {
	g, err := New(Config{Muted: true, Screen: "hidden"})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	packed, err := assets.Files.ReadFile(g.scene.Music)
	if err != nil {
		t.Fatal(err)
	}
	zipped, err := gzip.NewReader(bytes.NewReader(packed))
	if err != nil {
		t.Fatal(err)
	}
	pcm, err := io.ReadAll(zipped)
	zipped.Close()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := wav.DecodeWithSampleRate(48000, bytes.NewReader(pcm))
	if err != nil {
		t.Fatal(err)
	}
	frames := decoded.Length() / 4
	loopStart := g.scene.Descriptor.MusicLoopStartFrame
	if loopStart <= 0 || frames < 48000*160 || frames > 48000*180 {
		t.Fatal("hidden music is missing its complete introduction and musical cycle")
	}
	const tailFrames, repeatFrames = 16, 48
	if _, err = g.stream.Seek((frames-tailFrames)*8, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	seam := make([]byte, (tailFrames+repeatFrames)*8)
	if _, err = io.ReadFull(g.stream, seam); err != nil {
		t.Fatal("hidden music stopped instead of repeating:", err)
	}
	if _, err = decoded.Seek(loopStart*4, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	want := make([]byte, repeatFrames*4)
	if _, err = io.ReadFull(decoded, want); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < repeatFrames*2; i++ {
		expected := float32(int16(binary.LittleEndian.Uint16(want[i*2:]))) / 32768
		got := math.Float32frombits(binary.LittleEndian.Uint32(seam[tailFrames*8+i*4:]))
		if got != expected {
			t.Fatalf("loop sample %d repeats the wrong part of the recording: %g != %g", i, got, expected)
		}
	}
	if _, err = g.stream.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if _, err = decoded.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if _, err = io.ReadFull(decoded, want); err != nil {
		t.Fatal(err)
	}
	got := make([]byte, repeatFrames*8)
	if _, err = io.ReadFull(g.stream, got); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < repeatFrames*2; i++ {
		expected := float32(int16(binary.LittleEndian.Uint16(want[i*2:]))) / 32768
		if math.Float32frombits(binary.LittleEndian.Uint32(got[i*4:])) != expected {
			t.Fatal("restarting the hidden screen skipped its introduction")
		}
	}
}
