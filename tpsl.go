package agent2

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/varga-lp/data/klines"
)

type TPSL struct {
	TakeProfit float64 `json:"tp"`
	StopLoss   float64 `json:"sl"`
}

const (
	maxTPSL  = float64(0.03)  // %3.0
	minTPSL  = float64(0.004) // %0.4
	tpSLStep = float64(0.001) // %0.1
)

func randTreshold() float64 {
	r := rand.Float64()*(maxTPSL-minTPSL) + minTPSL
	m := (1.0 / tpSLStep)

	return math.Round(r*m) / m
}

var (
	ErrTresholdIsOutsideOfBoundries = fmt.Errorf("treshold is outside of boundries")
	ErrPositionCantBeNilForTP       = fmt.Errorf("position can't be nil for TP")
)

func randTresholdGTE(num float64) (float64, error) {
	if num > maxTPSL || num < minTPSL {
		return 0, ErrTresholdIsOutsideOfBoundries
	}

	for {
		r := randTreshold()

		if r >= num {
			return r, nil
		}
	}
}

func RandomTPSL() *TPSL {
	sl := randTreshold()
	tp, _ := randTresholdGTE(sl)

	return &TPSL{
		TakeProfit: tp,
		StopLoss:   sl,
	}
}

func (ts *TPSL) TPNetClose(pos *Position, closeLong klines.Kline, closeShort klines.Kline) (bool, error) {
	if pos == nil {
		return false, ErrPositionCantBeNilForTP
	}

	if (pos.NetProfit(closeLong, closeShort) / defaultAllocation) >= ts.TakeProfit {
		return true, nil
	}
	return false, nil
}

func (ts *TPSL) SLNetClose(pos *Position, closeLong klines.Kline, closeShort klines.Kline) (bool, error) {
	if pos == nil {
		return false, ErrPositionCantBeNilForTP
	}

	if -(pos.NetProfit(closeLong, closeShort) / defaultAllocation) >= ts.StopLoss {
		return true, nil
	}
	return false, nil
}
