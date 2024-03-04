package agent2

import (
	"math/rand"
	"testing"
	"time"
)

func TestRandBackOff(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		rb := randBackoff()

		if rb > maxBackoffMillis || rb < minBackoffMillis {
			t.Errorf("backoff is outside of boundries")
		}
	}
}

func TestRandBackOff_SingleIteration(t *testing.T) {
	rand.Seed(1)

	r := randBackoff()
	expected := int64(510000)

	if r != expected {
		t.Errorf("backoff %d is not expected %d", r, expected)
	}
}

func TestBackoff_TradeAllowed_NilLastTrade(t *testing.T) {
	bo := RandomBackoff()

	if !bo.TradeAllowed(nil) {
		t.Errorf("trade should be allowed when there is no last trade")
	}
}

func TestBackoff_TradeAllowed_CloseLastTrade(t *testing.T) {
	bo := RandomBackoff()
	bo.DurationMillis = 1_000

	tr := &Trade{
		CloseTime: time.Now().UnixMilli() - 500,
	}

	if bo.TradeAllowed(tr) {
		t.Errorf("trade should not be allowed when there is a close trade")
	}
}

func TestBackoff_TradeAllowed_FarawayLastTrade(t *testing.T) {
	bo := RandomBackoff()
	bo.DurationMillis = 1_000

	tr := &Trade{
		CloseTime: time.Now().UnixMilli() - 1001,
	}

	if !bo.TradeAllowed(tr) {
		t.Errorf("trade should be allowed when there is a far away trade")
	}
}
