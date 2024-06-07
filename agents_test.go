package agent2

import (
	"math/rand"
	"testing"
	"time"
)

func TestAgentMarshal_NoIndicators(t *testing.T) {
	rand.Seed(0)

	ag := &Agent{
		Tpsl:         RandomTPSL(),
		Backoff:      RandomBackoff(),
		ExpiryMillis: randExpiry(),
	}

	pload, err := ag.Marshal()
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	expected := `{"tpsl":{"tp":0.0145,"sl":0.0145},"backoff":{"mls":960000},"expiry_mls":18720000,"bbs":null,"rsis":null}`

	if string(pload) != expected {
		t.Errorf("expected payload %s, received %s", expected, string(pload))
	}
}

// func TestAgentMarshal_WithIndicators(t *testing.T) {
// 	rand.Seed(1)

// 	ag := RandomAgent()
// 	pload, _ := ag.Marshal()

// 	expected := `{"tpsl":{"tp":0.028,"sl":0.02},"backoff":{"mls":390000},"expiry_mls":16560000,"bbs":[{"mon":14,"val_pos":0,"line":2,"period":98,"multiplier":2.61},{"mon":8,"val_pos":1,"line":2,"period":105,"multiplier":2.1239},{"mon":11,"val_pos":1,"line":2,"period":146,"multiplier":1.8541}],"rsis":[{"mon":8,"val_pos":0,"target_val":69,"period":13},{"mon":2,"val_pos":1,"target_val":66,"period":55},{"mon":14,"val_pos":0,"target_val":69,"period":89},{"mon":9,"val_pos":0,"target_val":53,"period":113}]}`

// 	if string(pload) != expected {
// 		t.Errorf("expected payload %s, received %s", expected, string(pload))
// 	}
// }

func TestUnmarshalAgent_WithDummyPayload(t *testing.T) {
	pload := []byte("dummy")

	_, err := UnmarshalAgent(pload)
	if err == nil {
		t.Errorf("expected error but nothing raised")
	}
}

func TestUnmarshalAgent_Tpsl(t *testing.T) {
	ag := RandomAgent()
	pload, _ := ag.Marshal()

	ag2, _ := UnmarshalAgent(pload)
	if ag2.Tpsl.TakeProfit != ag.Tpsl.TakeProfit {
		t.Errorf("take profits does not match")
	}
	if ag2.Tpsl.StopLoss != ag.Tpsl.StopLoss {
		t.Errorf("stop losses does not match")
	}
}

func TestUnmarshalAgent_Backoff(t *testing.T) {
	ag := RandomAgent()
	pload, _ := ag.Marshal()

	ag2, _ := UnmarshalAgent(pload)
	if ag2.Backoff.DurationMillis != ag.Backoff.DurationMillis {
		t.Errorf("backoffs does not match")
	}
}

func TestUnmarshalAgent_ExpiryMillis(t *testing.T) {
	ag := RandomAgent()
	pload, _ := ag.Marshal()

	ag2, _ := UnmarshalAgent(pload)
	if ag2.ExpiryMillis != ag.ExpiryMillis {
		t.Errorf("expiry millis does not match")
	}
}

func TestUnmarshalAgent_IndicatorLens(t *testing.T) {
	ag := RandomAgent()
	pload, _ := ag.Marshal()

	ag2, _ := UnmarshalAgent(pload)
	if len(ag.Bbs) != len(ag2.Bbs) {
		t.Errorf("bb lens does not match")
	}
	if len(ag.Rsis) != len(ag2.Rsis) {
		t.Errorf("rsis lens does not match")
	}
}

func TestUnmarshalAgent_FirstBB(t *testing.T) {
	ag := RandomAgent()
	pload, _ := ag.Marshal()

	ag2, _ := UnmarshalAgent(pload)

	if ag2.Bbs[0].Mon != ag.Bbs[0].Mon {
		t.Errorf("first bb mons does not match")
	}
	if ag2.Bbs[0].ValuePos != ag.Bbs[0].ValuePos {
		t.Errorf("first bb value posses does not match")
	}
	if ag2.Bbs[0].Line != ag.Bbs[0].Line {
		t.Errorf("first bb lines does not match")
	}
	if ag2.Bbs[0].Period != ag.Bbs[0].Period {
		t.Errorf("first bb periods does not match")
	}
	if ag2.Bbs[0].Multiplier != ag.Bbs[0].Multiplier {
		t.Errorf("first bb multipliers does not match")
	}
}

func TestUnmarshalAgent_FirstRSI(t *testing.T) {
	ag := RandomAgent()
	pload, _ := ag.Marshal()

	ag2, _ := UnmarshalAgent(pload)

	if ag2.Rsis[0].Mon != ag.Rsis[0].Mon {
		t.Errorf("first rsi mons does not match")
	}
	if ag2.Rsis[0].ValuePos != ag.Rsis[0].ValuePos {
		t.Errorf("first rsi value posses does not match")
	}
	if ag2.Rsis[0].TargetVal != ag.Rsis[0].TargetVal {
		t.Errorf("first rsi target vals does not match")
	}
	if ag2.Rsis[0].Period != ag.Rsis[0].Period {
		t.Errorf("first rsi periods does not match")
	}
}

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
