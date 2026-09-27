package screens

import "github.com/olivierh59500/democonstructionkit/motion"

// Authored phrase poses follow the recording's alternating stacks, reversed
// travel, arches and diagonal paths. CuedFormation performs the interpolation.
func replicantsFormation() (*motion.CuedFormation, error) {
	pose := func(time, x, y, spacing float64) motion.FormationPoseKey {
		return motion.FormationPoseKey{
			Time: time, Origin: motion.Point{X: x, Y: y},
			Spacing: motion.Point{X: spacing}, Ease: motion.Linear,
		}
	}
	bend := func(time, x, y, spacingX, spacingY, arcY float64) motion.FormationPoseKey {
		return motion.FormationPoseKey{
			Time: time, Origin: motion.Point{X: x, Y: y},
			Spacing: motion.Point{X: spacingX, Y: spacingY},
			Arc:     motion.Point{Y: arcY}, Ease: motion.Linear,
		}
	}
	poses := []motion.FormationPoseKey{
		pose(0, 257, 237, 20),
		pose(.2, 184, 231, 30),
		pose(.4, 147, 237, 26),
		pose(.6, 127, 250, 21),
		pose(.8, 192, 254, 27.5),
		pose(1, 330, 250, 24.5),
		pose(1.2, 488, 242, 13.5),
	}
	// This two-second left/right sweep repeats before the letters open out.
	sweep := []motion.FormationPoseKey{
		pose(0, 664, 238, -9.3),
		pose(.2, 664, 237, -20.8),
		pose(.4, 601, 237, -27.5),
		pose(.6, 463, 237, -24.5),
		pose(.8, 305, 240, -13.5),
		pose(1, 129, 250, 9.3),
		pose(1.2, 127, 250, 21),
		pose(1.4, 192, 254, 27.5),
		pose(1.6, 330, 250, 24.5),
		pose(1.8, 488, 242, 13.5),
		pose(2, 664, 238, -9.3),
	}
	for cycle := 0; cycle < 3; cycle++ {
		for step, key := range sweep {
			if cycle > 0 && step == 0 || cycle == 2 && step > 5 {
				continue
			}
			key.Time += 1.4 + 2*float64(cycle)
			poses = append(poses, key)
		}
	}
	poses = append(poses,
		pose(6.6, 129, 250, 19.5),
		pose(6.8, 152, 237, 29),
		bend(7, 162, 232, 33.5, -2.7, 0),
		bend(8, 230, 184, 28.5, 2, -45),
		bend(9, 247, 173, 28, 4.4, -50),
		pose(10, 157, 231, 33.5),
		bend(11, 139, 190, 31.5, 10, 0),
		bend(12, 221, 199, 31.5, 9.3, 0),
		bend(13, 124, 230, 29, 3.5, -55),
		bend(14, 115, 234, 29.5, 2, -55),
		bend(15, 200, 210, 29, 9.5, 0),
		bend(16, 145, 234, 28.5, 7.5, 0),
		pose(17, 150, 248, 27),
		pose(18, 280, 238, 17),
		pose(19, 286, 254, 26.5),
		pose(20, 506, 240, -26.3),
		pose(21, 286, 254, 26.5),
		pose(22, 506, 240, -26.3),
		pose(23, 286, 254, 26.5),
		pose(24, 506, 240, -26.3),
		pose(25, 166, 237, 33.5),
		pose(25.8, 257, 237, 20),
		pose(26, 257, 237, 20),
	)
	return motion.NewCuedFormation(motion.CuedFormationConfig{
		Origin:  motion.Point{X: 162, Y: 230},
		Spacing: motion.Point{X: 33.5}, Count: 14, Loop: 26,
		PoseKeys: poses,
	})
}
