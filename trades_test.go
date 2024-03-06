package agent2

import "testing"

func TestNewTrade_NilPos(t *testing.T) {
	kln1, kln2 := dummyKlines(1)[0], dummyKlines(1)[0]

	if _, err := NewTrade(nil, NoReason, kln1, kln2); err != ErrPositionCantBeNilForTrade {
		t.Errorf("expected %v to be raised, raised %v", ErrPositionCantBeNilForTrade, err)
	}
}

func TestNewTrade_CloseTimeErr(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])

	longClose := dummyKlines(2)[1]
	longClose.CloseTime = dummyKlines(1)[0].CloseTime - 1

	if _, err := NewTrade(pos, NoReason, longClose, dummyKlines(2)[1]); err != ErrPositionCloseIsNotLTECloseTime {
		t.Errorf("expected %v to be raised, raised %v", ErrPositionCloseIsNotLTECloseTime, err)
	}
}

func TestNewTrade_OpenTime(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])

	tr, err := NewTrade(pos, NoReason, dummyKlines(2)[1], dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected nothing to be raised, raised %v", err)
	}

	if tr.OpenTime != dummyKlines(1)[0].OpenTime {
		t.Errorf("open time is not assigned properly")
	}
}

func TestNewTrade_CloseTime(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])

	tr, err := NewTrade(pos, NoReason, dummyKlines(2)[1], dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected nothing to be raised, raised %v", err)
	}

	if tr.CloseTime != dummyKlines(2)[1].CloseTime {
		t.Errorf("close time is not assigned properly")
	}
}

func TestNewTrade_ClosingReason(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])

	tr, err := NewTrade(pos, Expiry, dummyKlines(2)[1], dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected nothing to be raised, raised %v", err)
	}

	if tr.Reason != Expiry {
		t.Errorf("closing reason is not assigned properly")
	}
}

func TestNewTrade_NetProfit(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])

	tr, err := NewTrade(pos, NoReason, dummyKlines(2)[1], dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected nothing to be raised, raised %v", err)
	}

	expected := 248.9875
	if tr.NetProfit != expected {
		t.Errorf("net profit %.4f is not expected %.4f", tr.NetProfit, expected)
	}
}

func TestNewTrade_DurationSecs(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])
	closeLong := dummyKlines(2)[1]
	closeLong.CloseTime = 10000

	tr, err := NewTrade(pos, NoReason, closeLong, dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected nothing to be raised, raised %v", err)
	}

	expected := int64(10)
	if tr.DurationSecs != expected {
		t.Errorf("durationSecs %d is not expected %d", tr.DurationSecs, expected)
	}
}

func TestNewTrade_String(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])
	closeLong := dummyKlines(2)[1]
	closeLong.CloseTime = 10000

	tr, err := NewTrade(pos, NoReason, closeLong, dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected nothing to be raised, raised %v", err)
	}

	expected := "[trd] OT=0, CT=10000, DSecs=10, Net=248.99, R=noReason"
	if tr.String() != expected {
		t.Errorf("string %s is not expected %s", tr.String(), expected)
	}
}

func TestRoundTo4d(t *testing.T) {
	expected := 0.3333
	r := roundTo4d(float64(1.0) / float64(3.0))

	if r != expected {
		t.Errorf("expected float %.4f is not equal to r %.4f", expected, r)
	}
}

func TestNewBucket_InvalidArgs(t *testing.T) {
	if _, err := NewBucket(2, 1); err != ErrBucketEndTimeIsNotGTStartTime {
		t.Errorf("expected error not raised")
	}
}

func TestNewBucket_WithValidArgs(t *testing.T) {
	_, err := NewBucket(1, 2)
	if err != nil {
		t.Errorf("no error expected but raised")
	}
}

func TestBucket_LastTrade_EmptyTrades(t *testing.T) {
	b, _ := NewBucket(1, 2)

	if b.LastTrade() != nil {
		t.Errorf("expected nil trade")
	}
}

func TestBucket_LastTrade_WithTrades(t *testing.T) {
	pos1, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])
	pos2, _ := NewPosition(dummyKlines(2)[1], dummyKlines(2)[1])

	tr1, _ := NewTrade(pos1, NoReason, dummyKlines(2)[1], dummyKlines(2)[1])
	tr2, _ := NewTrade(pos2, NoReason, dummyKlines(3)[2], dummyKlines(3)[2])
	tr2.Reason = Expiry

	bu, _ := NewBucket(1, 2)
	bu.Trades = append(bu.Trades, tr1)
	bu.Trades = append(bu.Trades, tr2)

	if bu.LastTrade().Reason != Expiry {
		t.Errorf("unexpected reason")
	}
}

func TestBucket_AppendTrade(t *testing.T) {
	pos1, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])
	pos2, _ := NewPosition(dummyKlines(2)[1], dummyKlines(2)[1])

	bu, _ := NewBucket(1, 2)
	err := bu.AppendTrade(pos1, TakeProfit, dummyKlines(2)[1], dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected no error, but raised")
	}

	err = bu.AppendTrade(pos2, StopLoss, dummyKlines(3)[2], dummyKlines(3)[2])
	if err != nil {
		t.Errorf("expected no error, but raised")
	}

	if bu.LastTrade().Reason != StopLoss {
		t.Errorf("unexpected reason")
	}
}

func TestBucket_HitRatio_NoTrades(t *testing.T) {
	bu, _ := NewBucket(1, 2)
	if bu.HitRatio() != 0.0 {
		t.Errorf("unpexpected hit ratio")
	}
}

func TestBucket_HitRatio_WithTrades(t *testing.T) {
	bu, _ := NewBucket(1, 2)
	bu.Trades = append(bu.Trades, &Trade{NetProfit: -0.1})
	bu.Trades = append(bu.Trades, &Trade{NetProfit: 0.5})
	bu.Trades = append(bu.Trades, &Trade{NetProfit: 0.8})

	expected := 0.6666
	if bu.HitRatio() != expected {
		t.Errorf("unexpected hit ratio")
	}
}

func TestBucket_ProfitPerDay_NoTrades(t *testing.T) {
	bu, _ := NewBucket(1, 2)
	if bu.ProfitPerDay() != 0.0 {
		t.Errorf("unexpected profit per day")
	}
}

func TestBucket_ProfitPerDay_WithTrades(t *testing.T) {
	bu, _ := NewBucket(0, int64(dayLenMillis)*2)
	bu.Trades = append(bu.Trades, &Trade{NetProfit: 10.0})
	bu.Trades = append(bu.Trades, &Trade{NetProfit: 20.0})
	bu.Trades = append(bu.Trades, &Trade{NetProfit: -10.0})

	expected := 10.0000
	if bu.ProfitPerDay() != expected {
		t.Errorf("unexpected profit per day")
	}
}
