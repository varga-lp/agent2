package agent2

import (
	"math/rand"
	"testing"
)

func TestRandPeriod(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		p := randPeriod()

		if p < minPeriod || p > maxPeriod {
			t.Errorf("period %d is outside of allowed interval", p)
		}
	}
}

func TestRandPeriodSingleIteration(t *testing.T) {
	rand.Seed(0)

	expected, p := 244, randPeriod()
	if p != expected {
		t.Errorf("period %d is not expected %d", p, expected)
	}
}

func TestRandMultiplier(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		m := randMultiplier()

		if m < minMultiplier || m > maxMultiplier {
			t.Errorf("multiplier %.2f is outside of allowed interval", m)
		}
	}
}

func TestRandMultiplierSingleIteration(t *testing.T) {
	rand.Seed(0)

	expected, m := 4.7534, randMultiplier()
	if m != expected {
		t.Errorf("multiplier %.4f is not expected %.4f", m, expected)
	}
}
