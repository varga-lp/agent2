package agent2

import "testing"

func TestNewTrade_NilPos(t *testing.T) {
	kln1, kln2 := dummyKlines(1)[0], dummyKlines(1)[0]

	if _, err := NewTrade(nil, kln1, kln2); err != ErrPositionCantBeNil {
		t.Errorf("expected %v to be raised, raised %v", ErrPositionCantBeNil, err)
	}
}

func TestNewTrade_OpenTime(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])

	tr, err := NewTrade(pos, dummyKlines(2)[1], dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected nothing to be raised, raised %v", err)
	}

	if tr.OpenTime != dummyKlines(1)[0].OpenTime {
		t.Errorf("open time is not assigned properly")
	}
}

func TestNewTrade_CloseTime(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])

	tr, err := NewTrade(pos, dummyKlines(2)[1], dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected nothing to be raised, raised %v", err)
	}

	if tr.CloseTime != dummyKlines(2)[1].CloseTime {
		t.Errorf("close time is not assigned properly")
	}
}

func TestNewTrade_NetProfit(t *testing.T) {
	pos, _ := NewPosition(dummyKlines(1)[0], dummyKlines(1)[0])

	tr, err := NewTrade(pos, dummyKlines(2)[1], dummyKlines(2)[1])
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

	tr, err := NewTrade(pos, closeLong, dummyKlines(2)[1])
	if err != nil {
		t.Errorf("expected nothing to be raised, raised %v", err)
	}

	expected := int64(10)
	if tr.DurationSecs != expected {
		t.Errorf("durationSecs %d is not expected %d", tr.DurationSecs, expected)
	}
}
