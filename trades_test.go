package agent2

import "testing"

func TestNewTrade_NilPos(t *testing.T) {
	kln1, kln2 := dummyKlines(1)[0], dummyKlines(1)[0]

	if _, err := NewTrade(nil, NoReason, kln1, kln2); err != ErrPositionCantBeNilForTrade {
		t.Errorf("expected %v to be raised, raised %v", ErrPositionCantBeNilForTrade, err)
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

	expected := "[trd] OT=0, DSecs=10, Net=248.99, R=noReason"
	if tr.String() != expected {
		t.Errorf("string %s is not expected %s", tr.String(), expected)
	}
}
