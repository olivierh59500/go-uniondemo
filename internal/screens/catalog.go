// Package screens contains the individual Union compositions.
package screens

// Descriptor identifies a composition and its default native presentation.
type Descriptor struct {
	ID, Title, Directory, Music string
	Width, Height               int
}

var catalog = []Descriptor{
	{"beatdis", "TCB — Beat Dis", "beatdis", "audio/beatdis.ym", 768, 536},
	{"delta", "Delta Force", "delta", "audio/delta.ym", 768, 536},
	{"tnt3", "TNT Crew — Screen 3", "tnt3", "audio/tnt3.ym", 768, 536},
	{"wow", "TCB — Wow Scroller", "wow", "audio/wow.ym", 768, 536},
	{"hidden", "Hidden Screen", "hidden", "audio/hidden.ym", 768, 536},
	{"starballs", "TNT Crew — Starballs", "starballs", "audio/starballs.ym", 768, 536},
	{"replicants", "The Replicants", "replicants", "audio/replicants.ym", 768, 536},
	{"tnt2", "TNT Crew — Screen 2", "tnt2", "audio/tnt2.ym", 768, 536},
	{"level16", "Level 16", "level16", "audio/level16.ym", 768, 536},
	{"multiplane", "TCB — Multi-Plane 3D Scroller", "multiplane", "audio/multiplane.ym", 768, 536},
	{"diskcopier", "TEX — Disk Copier", "diskcopier", "audio/diskcopier.ym", 768, 536},
}

func Catalog() []Descriptor { return append([]Descriptor(nil), catalog...) }

func Find(id string) (Descriptor, bool) {
	for _, d := range catalog {
		if d.ID == id {
			return d, true
		}
	}
	return Descriptor{}, false
}
