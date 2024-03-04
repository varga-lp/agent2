package agent2

import (
	"testing"
	"time"
)

func TestRandomAgent_BBLen(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		bbLen := len(RandomAgent().Bbs)

		if bbLen < 1 || bbLen > maxBBCount {
			t.Errorf("bbLen is outside of boundries")
		}
	}
}

func TestRandomAgent_RSILen(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		rsiLen := len(RandomAgent().Rsis)

		if rsiLen < 1 || rsiLen > maxRSICount {
			t.Errorf("rsiLen is outside of boundries")
		}
	}
}

func TestRandomAgent_BBSorted(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		ag := RandomAgent()

		lastPeriod := 0
		for _, bb := range ag.Bbs {
			if bb.Period < lastPeriod {
				t.Errorf("bbs are not sorted according to period")
			}
			lastPeriod = bb.Period
		}
	}
}

func TestRandomAgent_RSISorted(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		ag := RandomAgent()

		lastPeriod := 0
		for _, rsi := range ag.Rsis {
			if rsi.Period < lastPeriod {
				t.Errorf("rsis are not sorted according to period")
			}
			lastPeriod = rsi.Period
		}
	}
}

func TestOpenPos_RaiseError_WithLowNumberOfKlines(t *testing.T) {
	ag := RandomAgent()

	_, err := ag.OpenPos(dummyKlines(minActivationKlineLength-1), dummyKlines(minActivationKlineLength-1), nil)
	if err != ErrKlinesAreBelowMinActivationKlineLength {
		t.Errorf("expected %v error but raised %v", ErrKlinesAreBelowMinActivationKlineLength, err)
	}
}

func TestOpenPos_RaiseNoError_WithValidArgs_NoTrade(t *testing.T) {
	ag := RandomAgent()

	_, err := ag.OpenPos(dummyKlines(250), dummyKlines(250), nil)
	if err != nil {
		t.Errorf("expected no error but received %v", err)
	}
}

func TestOpenPos_RaiseNoError_WithValidArgs_Trade(t *testing.T) {
	ag := RandomAgent()
	tr := &Trade{
		CloseTime: time.Now().UnixMilli(),
	}

	_, err := ag.OpenPos(dummyKlines(250), dummyKlines(250), tr)
	if err != nil {
		t.Errorf("expected no error but received %v", err)
	}
}

func TestOpenPos_SingleIndicator_Active(t *testing.T) {
	ag := RandomAgent()
	ag.Bbs = make([]*BB, 0)

	rsi := RandomRSI()
	rsi.TargetVal = 10.0
	rsi.ValuePos = Above
	rsi.Period = 250
	ag.Rsis = []*RSI{rsi}

	open, _ := ag.OpenPos(dummyKlines(250), dummyKlines(250), nil)
	if !open {
		t.Errorf("expected openpos to return true, returned false")
	}
}

func TestOpenPos_SingleIndicator_NonActive(t *testing.T) {
	ag := RandomAgent()
	ag.Bbs = make([]*BB, 0)

	rsi := RandomRSI()
	rsi.TargetVal = 10.0
	rsi.ValuePos = Below
	rsi.Period = 250
	ag.Rsis = []*RSI{rsi}

	open, _ := ag.OpenPos(dummyKlines(250), dummyKlines(250), nil)
	if open {
		t.Errorf("expected openpos to return false, returned true")
	}
}

func TestOpenPos_SingleIndicator_OneActiveOneInactive(t *testing.T) {
	ag := RandomAgent()

	bb := RandomBB()
	bb.Period = 250
	bb.Line = Upper
	bb.ValuePos = Above
	ag.Bbs = []*BB{bb}

	rsi := RandomRSI()
	rsi.TargetVal = 10.0
	rsi.ValuePos = Above
	rsi.Period = 250
	ag.Rsis = []*RSI{rsi}

	open, _ := ag.OpenPos(dummyKlines(250), dummyKlines(250), nil)
	if open {
		t.Errorf("expected openpos to return false, returned true")
	}
}

func TestClosePos_RaiseWithNilPos(t *testing.T) {
	ag := RandomAgent()
	kln1, kln2 := dummyKlines(1)[0], dummyKlines(1)[0]

	_, _, err := ag.ClosePos(nil, kln1, kln2)
	if err != ErrPositionCantBeNilForClose {
		t.Errorf("expected closepos to raise %v, returned %v",
			ErrPositionCantBeNilForClose, err)
	}
}

func TestClosePos_StopLoss(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln1O.Close = 1.0
	kln1C.Close = 1.0
	kln2O.Close = 1.0
	kln2C.Close = 1.5

	p, _ := NewPosition(kln1O, kln2O)
	ag := RandomAgent()

	clos, reason, _ := ag.ClosePos(p, kln1C, kln2C)
	if !clos {
		t.Errorf("expected close to be true")
	}
	if reason != StopLoss {
		t.Errorf("expected reason to be stop loss")
	}
}

func TestClosePos_TakeProfit(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln1O.Close = 1.0
	kln1C.Close = 1.0
	kln2O.Close = 1.0
	kln2C.Close = 0.9

	p, _ := NewPosition(kln1O, kln2O)
	ag := RandomAgent()

	clos, reason, _ := ag.ClosePos(p, kln1C, kln2C)
	if !clos {
		t.Errorf("expected close to be true")
	}
	if reason != TakeProfit {
		t.Errorf("expected reason to be take profit")
	}
}

func TestClosePos_Expiry(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln1O.Close = 1.0
	kln1C.Close = 1.0
	kln2O.Close = 1.0
	kln2C.Close = 1.0

	kln1O.CloseTime = time.Now().UnixMilli() - 10_000
	kln1C.CloseTime = time.Now().UnixMilli()

	p, _ := NewPosition(kln1O, kln2O)
	ag := RandomAgent()
	ag.ExpiryMillis = 9_999

	clos, reason, _ := ag.ClosePos(p, kln1C, kln2C)
	if !clos {
		t.Errorf("expected close to be true")
	}
	if reason != Expiry {
		t.Errorf("expected reason to be expiry")
	}
}
