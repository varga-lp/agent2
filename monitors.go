package agent2

import "math/rand"

type Monitor uint8

const (
	Close1 Monitor = iota
	Close2
	CloseR
	HighMLow1
	HighMLow2
	HighMLowR
	Volume1
	Volume2
	VolumeR
	TBVolOVol1
	TBVolOVol2
	TBVolOVolR
	Not1
	Not2
	NotR
)

var (
	primaryMons = []Monitor{
		CloseR, VolumeR, TBVolOVolR, NotR,
	}
	secondaryMons = []Monitor{
		Close1, Close2, HighMLow1, HighMLow2,
		HighMLowR, Volume1, Volume2, TBVolOVol1,
		TBVolOVol2, Not1, Not2,
	}
)

const (
	secondaryMonProb = 10
)

func randPrimaryMon() Monitor {
	return primaryMons[rand.Intn(len(primaryMons))]
}

func randSecondaryMon() Monitor {
	return secondaryMons[rand.Intn(len(secondaryMons))]
}

func randMon() Monitor {
	if rand.Intn(100) < secondaryMonProb {
		return randSecondaryMon()
	}
	return randPrimaryMon()
}
