package agent2

import (
	"math/rand"
	"testing"
)

func TestRandPrimaryMon(t *testing.T) {
	count := 0
	for _, mon := range primaryMons {
		for {
			if randPrimaryMon() == mon {
				count++
				break
			}
		}
	}

	if count != len(primaryMons) {
		t.Errorf("all primary mons not returned randomly")
	}
}

func TestRandSecondaryMon(t *testing.T) {
	count := 0
	for _, mon := range secondaryMons {
		for {
			if randSecondaryMon() == mon {
				count++
				break
			}
		}
	}

	if count != len(secondaryMons) {
		t.Errorf("all secondary mons not returned randomly")
	}
}

func TestRandMon(t *testing.T) {
	count := 0
	for _, mon := range primaryMons {
		for {
			if randMon() == mon {
				count++
				break
			}
		}
	}
	for _, mon := range secondaryMons {
		for {
			if randMon() == mon {
				count++
				break
			}
		}
	}

	if count != len(primaryMons)+len(secondaryMons) {
		t.Errorf("all mons not returned randomly")
	}
}

func TestRandMonSingleCall(t *testing.T) {
	rand.Seed(0)

	mon := randMon()
	expected := TBVolOVolR

	if mon != expected {
		t.Errorf("expected %v, received %v", expected, mon)
	}
}
