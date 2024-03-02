package agent2

import (
	"math"
	"math/rand"
	"testing"

	"github.com/varga-lp/data/klines"
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

func TestRandBB_AssignsBothValPos(t *testing.T) {
	count := 0
	for _, valPos := range []ValuePos{Above, Below} {
		for {
			bb := RandomBB()
			if bb.ValuePos == valPos {
				count++
				break
			}
		}
	}
	if count != 2 {
		t.Errorf("not all valPos assinged randomly")
	}
}

func TestRandBB_AssignsAllBBLines(t *testing.T) {
	count := 0
	for _, bbLine := range []BBLine{Lower, Middle, Upper} {
		for {
			bb := RandomBB()
			if bb.Line == bbLine {
				count++
				break
			}
		}
	}
	if count != 3 {
		t.Errorf("not all bbLines assinged randomly")
	}
}

func TestBB_Marshal(t *testing.T) {
	bb := RandomBB()

	b, err := bb.Marshal()
	if err != nil {
		t.Errorf("unexpected marshaling error %s", err.Error())
	}

	expected := 92
	if len(b) != expected {
		t.Errorf("expected %d, received %d", expected, len(b))
	}
}

func TestBB_UnMarshal(t *testing.T) {
	bb := RandomBB()
	pload, _ := bb.Marshal()

	bb2, err := UnMarshalBB(pload)
	if err != nil {
		t.Errorf("error %s not expected", err.Error())
	}

	if bb2.Mon != bb.Mon {
		t.Errorf("mons are not unmarshaled correctly %v, %v", bb2.Mon, bb.Mon)
	}
	if bb2.Multiplier != bb.Multiplier {
		t.Errorf("multipliers are not unmarshalled correctly %f, %f", bb2.Multiplier, bb.Multiplier)
	}
}

func TestKTMV_InvalidKlns1Len(t *testing.T) {
	klns1 := make([]klines.Kline, 10)
	klns2 := make([]klines.Kline, 20)

	if _, err := klinesToMonValues(Close1, 20, klns1, klns2); err == nil {
		t.Errorf("expected error nothing raised")
	} else {
		if err.Error() != "klns1 length 10 should be 20" {
			t.Errorf("unexpected error %v", err)
		}
	}
}

func TestKTMV_InvalidKlns2Len(t *testing.T) {
	klns1 := make([]klines.Kline, 20)
	klns2 := make([]klines.Kline, 10)

	if _, err := klinesToMonValues(Close1, 20, klns1, klns2); err == nil {
		t.Errorf("expected error nothing raised")
	} else {
		if err.Error() != "klns2 length 10 should be 20" {
			t.Errorf("unexpected error %v", err)
		}
	}
}

func TestKTMV_ReturnLenShouldBeEqualToPeriod(t *testing.T) {
	klns1, klns2 := make([]klines.Kline, 20), make([]klines.Kline, 20)

	res, err := klinesToMonValues(Close1, 20, klns1, klns2)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	} else if len(res) != 20 {
		t.Errorf("len of res is not 20 but %d", len(res))
	}
}

func dummyKlines(length int) []klines.Kline {
	res := make([]klines.Kline, length)

	for i := 0; i < length; i++ {
		res[i] = klines.Kline{
			Open:           float64(i),
			High:           float64(i + 2),
			Low:            float64(i),
			Close:          float64(i + 1),
			Volume:         float64(i * 1000),
			TakerBuyVolume: float64(i * 600),
			NumberOfTrades: int64(i * 10),
		}
	}
	return res
}

func TestKTMV_ReturnsExpectedClose1s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)

	res, _ := klinesToMonValues(Close1, 4, klns1, klns2)
	expected := []float64{1, 2, 3, 4}

	for i, exp := range expected {
		if res[i] != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedClose2s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(10)[6:]

	res, _ := klinesToMonValues(Close2, 4, klns1, klns2)
	expected := []float64{7, 8, 9, 10}

	for i, exp := range expected {
		if res[i] != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedCloseRs(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)

	res, _ := klinesToMonValues(CloseR, 4, klns1, klns2)
	expected := []float64{1, 1, 1, 1}

	for i, exp := range expected {
		if math.Round(res[i]*100.0)/100.0 != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedHighMLow1s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)

	res, _ := klinesToMonValues(HighMLow1, 4, klns1, klns2)
	expected := []float64{2, 2, 2, 2}

	for i, exp := range expected {
		if res[i] != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedHighMLow2s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)
	klns2[3].High = 10.0

	res, _ := klinesToMonValues(HighMLow2, 4, klns1, klns2)
	expected := []float64{2, 2, 2, 7}

	for i, exp := range expected {
		if res[i] != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedHighMLowRs(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)
	klns2[3].High = 15.0

	res, _ := klinesToMonValues(HighMLowR, 4, klns1, klns2)
	expected := []float64{1, 1, 1, 0.17}

	for i, exp := range expected {
		if math.Round(res[i]*100)/100.0 != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedVolume1s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)

	res, _ := klinesToMonValues(Volume1, 4, klns1, klns2)
	expected := []float64{0, 1000, 2000, 3000}

	for i, exp := range expected {
		if res[i] != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedVolume2s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)
	klns2[3].Volume = 356.0

	res, _ := klinesToMonValues(Volume2, 4, klns1, klns2)
	expected := []float64{0, 1000, 2000, 356}

	for i, exp := range expected {
		if res[i] != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedVolumeRs(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)

	res, _ := klinesToMonValues(VolumeR, 4, klns1, klns2)
	expected := []float64{0, 1, 1, 1}

	for i, exp := range expected {
		if math.Round(res[i]*100.0)/100.0 != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedTBVolOVol1s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)

	res, _ := klinesToMonValues(TBVolOVol1, 4, klns1, klns2)
	expected := []float64{0, 0.6, 0.6, 0.6}

	for i, exp := range expected {
		if math.Round(res[i]*100.0)/100.0 != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedTBVolOVol2s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)
	klns2[3].Volume = 10000.0

	res, _ := klinesToMonValues(TBVolOVol2, 4, klns1, klns2)
	expected := []float64{0, 0.6, 0.6, 0.18}

	for i, exp := range expected {
		if math.Round(res[i]*100.0)/100.0 != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedTBVolOVolRs(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)
	klns2[3].Volume = 10000.0

	res, _ := klinesToMonValues(TBVolOVolR, 4, klns1, klns2)
	expected := []float64{0, 1, 1, 3.33}

	for i, exp := range expected {
		if math.Round(res[i]*100.0)/100.0 != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedNot1s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)

	res, _ := klinesToMonValues(Not1, 4, klns1, klns2)
	expected := []float64{0, 10.0, 20.0, 30.0}

	for i, exp := range expected {
		if res[i] != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedNot2s(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)
	klns2[3].NumberOfTrades = 123

	res, _ := klinesToMonValues(Not2, 4, klns1, klns2)
	expected := []float64{0, 10.0, 20.0, 123.0}

	for i, exp := range expected {
		if res[i] != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}

func TestKTMV_ReturnsExpectedNotRs(t *testing.T) {
	klns1, klns2 := dummyKlines(4), dummyKlines(4)
	klns2[3].NumberOfTrades = 150

	res, _ := klinesToMonValues(NotR, 4, klns1, klns2)
	expected := []float64{0, 1, 1, 0.2}

	for i, exp := range expected {
		if math.Round(res[i]*100.0)/100.0 != exp {
			t.Errorf("i %d is %.2f but expected %.2f", i, res[i], exp)
		}
	}
}
