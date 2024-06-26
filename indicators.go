package agent2

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/varga-lp/data/klines"
)

type ValuePos uint8

const (
	Above ValuePos = iota
	Below
)

type BBLine uint8

const (
	Lower BBLine = iota
	Middle
	Upper
)

type BB struct {
	Mon        Monitor  `json:"mon"`
	ValuePos   ValuePos `json:"val_pos"`
	Line       BBLine   `json:"line"`
	Period     int      `json:"period"`
	Multiplier float64  `json:"multiplier"`
}

const (
	minPeriod     = 10
	maxPeriod     = 250
	minMultiplier = float64(0.5)
	maxMultiplier = float64(5.0)
)

func randPeriod() int {
	return rand.Intn(maxPeriod-minPeriod) + minPeriod
}

func randMultiplier() float64 {
	r := rand.Float64()*(maxMultiplier-minMultiplier) + minMultiplier

	return math.Round(r*10_000.0) / 10_000.0
}

func RandomBB() *BB {
	return &BB{
		Mon:        randMon(),
		ValuePos:   ValuePos(rand.Intn(2)),
		Line:       BBLine(rand.Intn(3)),
		Period:     randPeriod(),
		Multiplier: randMultiplier(),
	}
}

const (
	epsilon = float64(0.0000000001)
)

func klinesToMonValues(mon Monitor, period int, klns1 []klines.Kline,
	klns2 []klines.Kline) ([]float64, error) {
	if len(klns1) != period {
		return nil, fmt.Errorf("klns1 length %d should be %d", len(klns1), period)
	}
	if len(klns2) != period {
		return nil, fmt.Errorf("klns2 length %d should be %d", len(klns2), period)
	}

	res := make([]float64, len(klns1))
	for i := 0; i < len(res); i++ {
		switch mon {
		case Close1:
			res[i] = klns1[i].Close
		case Close2:
			res[i] = klns2[i].Close
		case CloseR:
			res[i] = klns1[i].Close / (klns2[i].Close + epsilon)
		case HighMLow1:
			res[i] = klns1[i].High - klns1[i].Low
		case HighMLow2:
			res[i] = klns2[i].High - klns2[i].Low
		case HighMLowR:
			hml2 := klns2[i].High - klns2[i].Low
			if hml2 == 0 {
				res[i] = 0
			} else {
				res[i] = (klns1[i].High - klns1[i].Low) / hml2
			}
		case Volume1:
			res[i] = klns1[i].Volume
		case Volume2:
			res[i] = klns2[i].Volume
		case VolumeR:
			res[i] = klns1[i].Volume / (klns2[i].Volume + epsilon)
		case TBVolOVol1:
			res[i] = klns1[i].TakerBuyVolume / (klns1[i].Volume + epsilon)
		case TBVolOVol2:
			res[i] = klns2[i].TakerBuyVolume / (klns2[i].Volume + epsilon)
		case TBVolOVolR:
			t1 := klns1[i].TakerBuyVolume / (klns1[i].Volume + epsilon)
			t2 := klns2[i].TakerBuyVolume / (klns2[i].Volume + epsilon)
			res[i] = t1 / (t2 + epsilon)
		case Not1:
			res[i] = float64(klns1[i].NumberOfTrades)
		case Not2:
			res[i] = float64(klns2[i].NumberOfTrades)
		case NotR:
			res[i] = float64(klns1[i].NumberOfTrades) / (float64(klns2[i].NumberOfTrades) + epsilon)
		default:
			return nil, fmt.Errorf("mon %d is not defined", mon)
		}
	}
	return res, nil
}

func mean(vals []float64) (float64, error) {
	vlen := len(vals)
	if vlen == 0 {
		return 0, fmt.Errorf("vals len 0, can't take mean")
	}

	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals)), nil
}

func stddev(vals []float64, mean float64) (float64, error) {
	vlen := len(vals)
	if vlen == 0 {
		return 0, fmt.Errorf("vals len 0, can't take stddev")
	}

	sumOfSquaredDiffs := 0.0
	for _, v := range vals {
		sumOfSquaredDiffs += (v - mean) * (v - mean)
	}
	variance := sumOfSquaredDiffs / float64(len(vals))
	return math.Sqrt(variance), nil
}

func (bb *BB) Active(klns1 []klines.Kline, klns2 []klines.Kline) (bool, error) {
	vals, err := klinesToMonValues(bb.Mon, bb.Period, klns1, klns2)
	if err != nil {
		return false, err
	}

	mn, err := mean(vals)
	if err != nil {
		return false, err
	}
	std, err := stddev(vals, mn)
	if err != nil {
		return false, err
	}

	lastVal := vals[len(vals)-1]
	switch bb.Line {
	case Lower:
		if bb.ValuePos == Above {
			return lastVal > (mn - bb.Multiplier*std), nil
		} else if bb.ValuePos == Below {
			return lastVal < (mn - bb.Multiplier*std), nil
		}
	case Middle:
		if bb.ValuePos == Above {
			return lastVal > mn, nil
		} else if bb.ValuePos == Below {
			return lastVal < mn, nil
		}
	case Upper:
		if bb.ValuePos == Above {
			return lastVal > (mn + bb.Multiplier*std), nil
		} else if bb.ValuePos == Below {
			return lastVal < (mn + bb.Multiplier*std), nil
		}
	default:
		return false, fmt.Errorf("bbline %d is not defined", bb.Line)
	}
	return false, nil
}

const (
	minTVal = 5
	maxTVal = 95
)

type RSI struct {
	Mon       Monitor  `json:"mon"`
	ValuePos  ValuePos `json:"val_pos"`
	TargetVal float64  `json:"target_val"`
	Period    int      `json:"period"`
}

func randTargetVal() float64 {
	tval := minTVal + rand.Intn(maxTVal-minTVal)

	return float64(tval)
}

func RandomRSI() *RSI {
	return &RSI{
		Mon:       randMon(),
		ValuePos:  ValuePos(rand.Intn(2)),
		TargetVal: randTargetVal(),
		Period:    randPeriod(),
	}
}

func calcRsi(values []float64) (float64, error) {
	if len(values) < 2 {
		return 0, fmt.Errorf("needs min 2 elements to calculate rsi")
	}

	var gains, losses float64
	for i := 1; i < len(values); i++ {
		diff := values[i] - values[i-1]
		if diff > 0 {
			gains += diff
		} else {
			losses -= diff
		}
	}
	if losses == 0 {
		return 100, nil
	}
	return 100 - (100 / (1 + (gains / losses))), nil
}

func (rsi *RSI) Active(klns1 []klines.Kline, klns2 []klines.Kline) (bool, error) {
	vals, err := klinesToMonValues(rsi.Mon, rsi.Period, klns1, klns2)
	if err != nil {
		return false, err
	}

	r, err := calcRsi(vals)
	if err != nil {
		return false, err
	}

	switch rsi.ValuePos {
	case Above:
		return r > rsi.TargetVal, nil
	case Below:
		return r < rsi.TargetVal, nil
	default:
		return false, fmt.Errorf("valuePos %v is not defined", rsi.ValuePos)
	}
}
