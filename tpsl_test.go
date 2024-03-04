package agent2

import (
	"math/rand"
	"testing"
)

func TestRandTreshold(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		r := randTreshold()

		if r < minTPSL || r > maxTPSL {
			t.Errorf("randTreshold is outside of boundries")
		}
	}
}

func TestRandTreshold_SingleIteration(t *testing.T) {
	rand.Seed(0)

	r := randTreshold()
	expected := 0.0290

	if r != expected {
		t.Errorf("treshold %.4f is not equal to expected treshold %.4f", r, expected)
	}
}

func TestRandTresholdGTE_OutsideOfBoundries(t *testing.T) {
	if _, err := randTresholdGTE(maxTPSL + 0.0001); err != ErrTresholdIsOutsideOfBoundries {
		t.Errorf("expected error %v, raised %v", ErrTresholdIsOutsideOfBoundries, err)
	}
}

func TestRandTresholdGTE_ReturnsExpectedTreshold(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		r := randTreshold()

		r2, _ := randTresholdGTE(r)
		if r > r2 {
			t.Errorf("r %.4f is greater than r2 %.4f", r, r2)
		}
	}
}

func TestRandomTPSL_SingleIteration(t *testing.T) {
	rand.Seed(1)

	tpsl := RandomTPSL()

	if tpsl.TakeProfit != 0.0280 {
		t.Errorf("unexpected tp %.4f", tpsl.TakeProfit)
	}
	if tpsl.StopLoss != 0.0200 {
		t.Errorf("unexpected sl %.4f", tpsl.StopLoss)
	}
}

func TestTPNetClose_NilPos(t *testing.T) {
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	if _, err := RandomTPSL().TPNetClose(nil, kln1C, kln2C); err != ErrPositionCantBeNilForTP {
		t.Errorf("expected error %v but raised %v", ErrPositionCantBeNilForTP, err)
	}
}

func TestTPNetClose_NoClose(t *testing.T) {
	tpsl := RandomTPSL()
	tpsl.TakeProfit = 0.001

	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	pos, _ := NewPosition(kln1O, kln2O)

	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	clos, _ := tpsl.TPNetClose(pos, kln1C, kln2C)
	if clos {
		t.Errorf("unxpected close")
	}
}

func TestTPNetClose_Close(t *testing.T) {
	tpsl := RandomTPSL()
	tpsl.TakeProfit = 0.01

	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1O.Close = 100.0
	pos, _ := NewPosition(kln1O, kln2O)

	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C.Close = 102.2

	clos, _ := tpsl.TPNetClose(pos, kln1C, kln2C)
	if !clos {
		t.Errorf("unxpected close")
	}
}

func TestSLNetClose_NilPos(t *testing.T) {
	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	if _, err := RandomTPSL().SLNetClose(nil, kln1C, kln2C); err != ErrPositionCantBeNilForTP {
		t.Errorf("expected error %v but raised %v", ErrPositionCantBeNilForTP, err)
	}
}

func TestSLNetClose_NoClose(t *testing.T) {
	tpsl := RandomTPSL()
	tpsl.StopLoss = 0.001

	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	pos, _ := NewPosition(kln1O, kln2O)

	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]

	clos, _ := tpsl.SLNetClose(pos, kln1C, kln2C)
	if clos {
		t.Errorf("unxpected close")
	}
}

func TestSLNetClose_Close(t *testing.T) {
	tpsl := RandomTPSL()
	tpsl.StopLoss = 0.01

	kln1O, kln2O := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1O.Close = 100.0
	pos, _ := NewPosition(kln1O, kln2O)

	kln1C, kln2C := dummyKlines(1)[0], dummyKlines(1)[0]
	kln1C.Close = 98.1

	clos, _ := tpsl.SLNetClose(pos, kln1C, kln2C)
	if !clos {
		t.Errorf("unxpected close")
	}
}
