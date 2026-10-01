package screens

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/go-uniondemo/assets"
)

// The native 68000 tables supply one entry and a repeating formation. DCK owns
// position storage and sampling; the adapter only selects the native time bank.
type replicantsMotion struct {
	entry, cycle               *motion.KeyframedFormation
	entrySeconds, cycleSeconds float64
}

func (r *replicantsMotion) At(seconds float64, index int) motion.Point {
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		return motion.Point{}
	}
	tick := int64(math.Floor(math.Max(0, seconds)*50 + 1e-9))
	if tick < 31 {
		return r.entry.At(float64(tick)/50, index)
	}
	return r.cycle.At(float64((tick-31)%1273)/50, index)
}

func replicantsFormation() (*replicantsMotion, error) {
	raw, err := assets.Files.ReadFile("replicants/native-motion.json")
	if err != nil {
		return nil, fmt.Errorf("replicants motion: %w", err)
	}
	var source struct {
		Rate, Count  int
		Entry, Cycle [][][2]float64
	}
	if err := json.Unmarshal(raw, &source); err != nil {
		return nil, err
	}
	if source.Rate != 50 || source.Count != 13 || len(source.Entry) != 31 || len(source.Cycle) != 1273 {
		return nil, fmt.Errorf("replicants: incomplete native trajectory bank")
	}
	compile := func(bank [][][2]float64, loop bool) (*motion.KeyframedFormation, error) {
		frames := make([]motion.FormationFrame, len(bank)+1)
		for i, row := range bank {
			if len(row) != source.Count {
				return nil, fmt.Errorf("replicants: invalid pose %d", i)
			}
			points := make([]motion.Point, len(row))
			for j, p := range row {
				points[j] = motion.Point{X: p[0], Y: p[1]}
			}
			// The Atari changes positions once per VBL, without interpolation.
			frames[i] = motion.FormationFrame{Time: float64(i) / float64(source.Rate), Points: points, Ease: func(float64) float64 { return 0 }}
		}
		closing := frames[len(bank)-1].Points
		if loop {
			closing = frames[0].Points
		}
		frames[len(bank)] = motion.FormationFrame{Time: float64(len(bank)) / float64(source.Rate), Points: closing}
		duration := 0.0
		if loop {
			duration = frames[len(bank)].Time
		}
		return motion.NewKeyframedFormation(motion.KeyframedFormationConfig{Count: source.Count, Loop: duration, Frames: frames})
	}
	entry, err := compile(source.Entry, false)
	if err != nil {
		return nil, err
	}
	cycle, err := compile(source.Cycle, true)
	if err != nil {
		return nil, err
	}
	return &replicantsMotion{entry, cycle, float64(len(source.Entry)) / 50, float64(len(source.Cycle)) / 50}, nil
}
