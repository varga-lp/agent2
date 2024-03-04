package agent2

import (
	"math/rand"
	"testing"
)

func TestRandExpiry(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		re := randExpiry()

		if re > maxExpiryMillis || re < minExpiryMillis {
			t.Errorf("expiry is outside of boundries")
		}
	}
}

func TestRandExpiry_SingleIteration(t *testing.T) {
	rand.Seed(1)

	r := randExpiry()
	expected := int64(15960000)

	if r != expected {
		t.Errorf("expiry %d is not expected %d", r, expected)
	}
}
