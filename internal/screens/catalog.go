// Package screens contains the individual Union compositions.
package screens

import "time"

// IntroDuration is one complete performance of the opening YM soundtrack.
const IntroDuration = 38420 * time.Millisecond

// Descriptor identifies a composition and its default native presentation.
type Descriptor struct {
	ID, Title, Directory, Music string
	Width, Height               int
	// MusicLoopStartFrame skips a recorded introduction on repeats at 48000 Hz.
	MusicLoopStartFrame int64
}

var catalog = []Descriptor{
	{"beatdis", "TCB — Beat Dis", "beatdis", "audio/beatdis.ym", 768, 536, 0},
	{"delta", "Delta Force", "delta", "audio/delta.ym", 768, 536, 0},
	{"tnt3", "TNT Crew — Screen 3", "tnt3", "audio/tnt3.ym", 768, 536, 0},
	{"wow", "TCB — Wow Scroller", "wow", "audio/wow.ym", 768, 536, 0},
	{"hidden", "Hidden Screen", "hidden", "audio/thalion-forever.wav.gz", 768, 536, 37044},
	{"starballs", "TNT Crew — Starballs", "starballs", "audio/starballs.ym", 768, 536, 0},
	{"replicants", "The Replicants", "replicants", "audio/replicants.ym", 768, 536, 0},
	{"tnt2", "TNT Crew — Screen 2", "tnt2", "audio/tnt2.ym", 768, 536, 0},
	{"level16", "Level 16", "level16", "audio/level16.ym", 768, 536, 0},
	{"multiplane", "TCB — Multi-Plane 3D Scroller", "multiplane", "audio/multiplane.ym", 768, 536, 0},
	{"diskcopier", "TEX — Disk Copier", "diskcopier", "audio/diskcopier.ym", 768, 536, 0},
}

func Catalog() []Descriptor { return append([]Descriptor(nil), catalog...) }

// Introduction is separate from the eleven doors in the hall.
func Introduction() Descriptor {
	return Descriptor{ID: "intro", Title: "The Union — Introduction", Directory: "intro", Music: "audio/intro.ym", Width: 768, Height: 536}
}

func Find(id string) (Descriptor, bool) {
	if id == "intro" {
		return Introduction(), true
	}
	for _, d := range catalog {
		if d.ID == id {
			return d, true
		}
	}
	return Descriptor{}, false
}
