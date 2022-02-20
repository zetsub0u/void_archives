package store

import "github.com/zetsub0u/void_archives/archive"

func LoadDummyData(store archive.Archive) error {
	ref := archive.Ref{
		URL:     "https://www.youtube.com/watch?v=DyRGHCfhbQI",
		Creator: "zetsubou",
		Runs: []archive.Run{{
			Boss:  "SSS Jizo",
			Mode:  archive.Memorial,
			Score: 47568,
			Team: archive.Team{
				Valk1: archive.Valk{
					Name:    "BKE",
					Rank:    "S2",
					Comment: "",
				},
				Valk2: archive.Valk{
					Name:    "HOS",
					Rank:    "SS",
					Comment: "clone pop off",
				},
				Valk3: archive.Valk{
					Name:    "SA",
					Rank:    "SSS",
					Comment: "",
				},
				Elf: &archive.Elf{
					Name: "T0",
					Rank: archive.Elf3Star,
				},
				Keys:   nil,
				Emblem: nil,
			},
		}},
		Parsed:   false,
		Verified: false,
	}

	return store.AddRef(ref)
}
