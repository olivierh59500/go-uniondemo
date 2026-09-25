package screens

import "github.com/olivierh59500/democonstructionkit/motion"

// The resting letters are one sprite width apart. Separate cues send a motion
// across the phrase from either end, reproducing the changing paths visible in
// the original machine's wobbly-sprites screen.
func replicantsFormation() (*motion.CuedFormation, error) {
	return motion.NewCuedFormation(motion.CuedFormationConfig{
		Origin:  motion.Point{X: 162, Y: 230},
		Spacing: motion.Point{X: 31.5},
		Count:   14,
		Loop:    25,
		Cues: []motion.FormationCue{
			{
				Start: 0, Duration: 2.3, LeadIndex: 13, Stagger: .045, Fade: .08,
				X: []motion.FormationHarmonic{
					{FirstAmplitude: 0, LastAmplitude: 35, Cycles: .5},
					{FirstAmplitude: -85, LastAmplitude: 20, Cycles: 1},
				},
			},
			{
				Start: 2.4, Duration: 1, LeadIndex: 13, Stagger: .1, Fade: .04,
				Y: []motion.FormationHarmonic{{FirstAmplitude: 65, LastAmplitude: 65, Cycles: 1}},
			},
			{
				Start: 4, Duration: 1.3, LeadIndex: 0, Stagger: .05, Fade: .06,
				Y: []motion.FormationHarmonic{{FirstAmplitude: 60, LastAmplitude: 0, Cycles: 1.5}},
			},
			{
				Start: 5.5, Duration: 3, LeadIndex: 13, Stagger: .07, Fade: .12,
				X: []motion.FormationHarmonic{{FirstAmplitude: -50, LastAmplitude: 90, Cycles: .5}},
				Y: []motion.FormationHarmonic{
					{FirstAmplitude: -65, LastAmplitude: 65, Cycles: 1},
					{FirstAmplitude: 35, LastAmplitude: 35, Cycles: .5},
				},
			},
			{
				Start: 9, Duration: 3, LeadIndex: 0, Stagger: .06, Fade: .12,
				X: []motion.FormationHarmonic{{FirstAmplitude: 90, LastAmplitude: -90, Cycles: 1}},
				Y: []motion.FormationHarmonic{{FirstAmplitude: 70, LastAmplitude: 70, Cycles: .5}},
			},
			{
				Start: 12.5, Duration: 2.7, LeadIndex: 13, Fade: .12,
				X: []motion.FormationHarmonic{
					{FirstAmplitude: 160, LastAmplitude: -160, Cycles: .5},
					{FirstAmplitude: -40, LastAmplitude: 40, Cycles: 1},
				},
				Y: []motion.FormationHarmonic{{FirstAmplitude: -35, LastAmplitude: 35, Cycles: 1}},
			},
			{
				Start: 15, Duration: 3.5, LeadIndex: 13, Stagger: .07, Fade: .15,
				X: []motion.FormationHarmonic{{FirstAmplitude: 90, LastAmplitude: 90, Cycles: 1, IndexPhase: .12}},
				Y: []motion.FormationHarmonic{{FirstAmplitude: 45, LastAmplitude: 45, Cycles: 1, Phase: 1.57, IndexPhase: .12}},
			},
			{
				Start: 19, Duration: 3, LeadIndex: 0, Stagger: .07, Fade: .12,
				X: []motion.FormationHarmonic{{FirstAmplitude: -100, LastAmplitude: 100, Cycles: 1}},
				Y: []motion.FormationHarmonic{{FirstAmplitude: 80, LastAmplitude: -80, Cycles: 1.5}},
			},
			{
				Start: 22.3, Duration: 1.7, LeadIndex: 13, Stagger: .05, Fade: .08,
				X: []motion.FormationHarmonic{{FirstAmplitude: 80, LastAmplitude: -80, Cycles: .5}},
				Y: []motion.FormationHarmonic{{FirstAmplitude: -60, LastAmplitude: 60, Cycles: 1}},
			},
		},
	})
}
