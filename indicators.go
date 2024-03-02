package agent2

import (
	"bytes"
	"encoding/gob"
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

type ClosingReason uint8

const (
	TakeProfit ClosingReason = iota
	StopLoss
	Expiry
)

type BB struct {
	Mon        Monitor
	ValuePos   ValuePos
	Line       BBLine
	Period     int
	Multiplier float64
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

func (bb *BB) Marshal() ([]byte, error) {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(bb); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func UnMarshalBB(pload []byte) (*BB, error) {
	res := &BB{}

	decoder := gob.NewDecoder(bytes.NewReader(pload))
	if err := decoder.Decode(res); err != nil {
		return nil, err
	}
	return res, nil
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
