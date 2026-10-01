package screens

import (
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"io"
	"math"
	"os"
	"testing"

	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/go-uniondemo/assets"
)

func TestReplicantsUsesEveryNativePositionAt50Hz(t *testing.T) {
	formation, err := replicantsFormation()
	if err != nil {
		t.Fatal(err)
	}
	data, err := assets.Files.ReadFile("replicants/native-motion.json")
	if err != nil {
		t.Fatal(err)
	}
	var source struct{ Entry, Cycle [][][2]float64 }
	if err := json.Unmarshal(data, &source); err != nil {
		t.Fatal(err)
	}
	for tick := 0; tick < 31+2*1273; tick++ {
		row := source.Entry[min(tick, 30)]
		if tick >= 31 {
			row = source.Cycle[(tick-31)%1273]
		}
		for i, p := range row {
			want := motion.Point{X: p[0], Y: p[1]}
			got := formation.At(float64(tick)/50, i)
			if got != want {
				t.Fatalf("native tick %d letter %d: got %+v want %+v", tick, i, got, want)
			}
		}
	}
	if p := formation.At(0, 0); p.X != 384 || p.Y != 260 {
		t.Fatalf("unexpected native opening stack %+v", p)
	}
	if formation.At(.62, 0) != formation.At(.62+25.46, 0) {
		t.Fatal("native cycle changed at its repeat")
	}
	if formation.At(0, 13) != (motion.Point{}) {
		t.Fatal("native formation must contain thirteen visible letters")
	}
}

func TestReplicantsDisplayRateDoesNotAccelerateNativeClock(t *testing.T) {
	formation, err := replicantsFormation()
	if err != nil {
		t.Fatal(err)
	}
	for frame := 0; frame < 60*55; frame++ {
		time := float64(frame) / 60
		nativeTime := math.Floor(time*50+1e-9) / 50
		for i := 0; i < 13; i++ {
			if formation.At(time, i) != formation.At(nativeTime, i) {
				t.Fatalf("60 Hz host altered native pose at frame %d letter %d", frame, i)
			}
		}
	}
}

// The fixture is captured by executing the original 68000 drawing routine.
func TestReplicantsMatchesNativeDrawingOffsets(t *testing.T) {
	formation, err := replicantsFormation()
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("testdata/replicants-native-offsets.bin.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 33501*6 {
		t.Fatal("incomplete native sprite fixture")
	}
	for offset := 0; offset < len(data); offset += 6 {
		tick := int(binary.BigEndian.Uint16(data[offset:]))
		index := int(binary.BigEndian.Uint16(data[offset+2:]))
		want := int(binary.BigEndian.Uint16(data[offset+4:]))
		p := formation.At(float64(tick)/50, index)
		x, y := int(p.X-64)/2, int(p.Y-60)/2
		if got := y*160 + (x/16)*8; got != want {
			t.Fatalf("native tick %d letter %d: offset %d want %d", tick, index, got, want)
		}
	}
}
