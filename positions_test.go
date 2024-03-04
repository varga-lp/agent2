package agent2

import (
	"math"
	"testing"
)

func TestNewPosition_NotEqualOpenTimes(t *testing.T) {
	kln1, kln2 := dummyKlines(1)[0], dummyKlines(1)[0]
	kln2.OpenTime = 1000

	if _, err := NewPosition(kln1, kln2); err != ErrLongShortOpenTimeNotEqual {
		t.Errorf("expected error %v but nothing raised", ErrLongShortOpenTimeNotEqual)
	}
}

func TestGrossProfit_BothNeutral(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	p, err := NewPosition(kln1O, kln2O)
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	profit := math.Round(p.GrossProfit(kln1C, kln2C)*10_000) / 10_000.0
	expected := 0.0

	if profit != expected {
		t.Errorf("expected %.4f profit but returned %.4f profit", expected, profit)
	}
}

func TestGrossProfit_LongProfitShortNeutral(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln1O.Close = 1.0
	kln1C.Close = 2.0

	p, err := NewPosition(kln1O, kln2O)
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	profit := math.Round(p.GrossProfit(kln1C, kln2C)*10_000) / 10_000.0
	expected := 500.0000

	if profit != expected {
		t.Errorf("expected %.4f profit but returned %.4f profit", expected, profit)
	}
}

func TestGrossProfit_LongNeutralShortProfit(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln2O.Close = 2.0
	kln2C.Close = 1.0

	p, err := NewPosition(kln1O, kln2O)
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	profit := math.Round(p.GrossProfit(kln1C, kln2C)*10_000) / 10_000.0
	expected := 500.0000

	if profit != expected {
		t.Errorf("expected %.4f profit but returned %.4f profit", expected, profit)
	}
}

func TestGrossProfit_LongLossSortNeutral(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln1O.Close = 1.0
	kln1C.Close = 0.8

	p, err := NewPosition(kln1O, kln2O)
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	profit := math.Round(p.GrossProfit(kln1C, kln2C)*10_000) / 10_000.0
	expected := -100.0000

	if profit != expected {
		t.Errorf("expected %.4f profit but returned %.4f profit", expected, profit)
	}
}

func TestGrossProfit_LongNeutralShortLoss(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln2O.Close = 2.0
	kln2C.Close = 4.0

	p, err := NewPosition(kln1O, kln2O)
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	profit := math.Round(p.GrossProfit(kln1C, kln2C)*10_000) / 10_000.0
	expected := -250.0000

	if profit != expected {
		t.Errorf("expected %.4f profit but returned %.4f profit", expected, profit)
	}
}

func TestGrossProfit_LongLossShortProfit(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln1O.Close = 1.0
	kln1C.Close = 0.8
	kln2O.Close = 1.0
	kln2C.Close = 0.8

	p, err := NewPosition(kln1O, kln2O)
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	profit := math.Round(p.GrossProfit(kln1C, kln2C)*10_000) / 10_000.0
	expected := 25.0000

	if profit != expected {
		t.Errorf("expected %.4f profit but returned %.4f profit", expected, profit)
	}
}

func TestNetProfit_LongNeutralShortNeutral(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	p, err := NewPosition(kln1O, kln2O)
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	profit := math.Round(p.NetProfit(kln1C, kln2C)*10_000) / 10_000.0
	expected := -0.9000

	if profit != expected {
		t.Errorf("expected %.4f profit but returned %.4f profit", expected, profit)
	}

}

func TestNetProfit_LongProfitShortNeutral(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln1O.Close = 1.0
	kln1C.Close = 1.2

	p, err := NewPosition(kln1O, kln2O)
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	profit := math.Round(p.NetProfit(kln1C, kln2C)*10_000) / 10_000.0
	expected := 99.0550

	if profit != expected {
		t.Errorf("expected %.4f profit but returned %.4f profit", expected, profit)
	}

}

func TestNetProfit_LongProfitShortLoss(t *testing.T) {
	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	kln1O.Close = 1.0
	kln1C.Close = 2.0
	kln2O.Close = 1.0
	kln2C.Close = 1.5

	p, err := NewPosition(kln1O, kln2O)
	if err != nil {
		t.Errorf("expected no error but raised %v", err)
	}

	profit := math.Round(p.NetProfit(kln1C, kln2C)*10_000) / 10_000.0
	expected := 332.2833

	if profit != expected {
		t.Errorf("expected %.4f profit but returned %.4f profit", expected, profit)
	}
}
