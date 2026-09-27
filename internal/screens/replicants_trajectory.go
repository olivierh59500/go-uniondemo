package screens

import (
	"encoding/json"
	"fmt"

	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/go-uniondemo/assets"
)

// The positions are authored from the Atari performance. DCK owns their
// per-letter interpolation and seamless loop; this screen only supplies data.
func replicantsFormation() (*motion.KeyframedFormation, error) {
	data, err := assets.Files.ReadFile("replicants/motion.json")
	if err != nil {
		return nil, fmt.Errorf("replicants motion: %w", err)
	}
	var source struct {
		Loop   float64 `json:"loop"`
		Frames []struct {
			Time   float64      `json:"time"`
			Points [][2]float64 `json:"points"`
		} `json:"frames"`
	}
	if err := json.Unmarshal(data, &source); err != nil {
		return nil, fmt.Errorf("replicants motion: %w", err)
	}
	frames := make([]motion.FormationFrame, len(source.Frames))
	for i, frame := range source.Frames {
		points := make([]motion.Point, len(frame.Points))
		for j, xy := range frame.Points {
			points[j] = motion.Point{X: xy[0], Y: xy[1]}
		}
		frames[i] = motion.FormationFrame{Time: frame.Time, Points: points}
	}
	return motion.NewKeyframedFormation(motion.KeyframedFormationConfig{
		Count: 14, Loop: source.Loop, Frames: frames,
	})
}
