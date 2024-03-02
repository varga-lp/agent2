package agent2

import (
	"math"
	"math/rand"

	"github.com/varga-lp/data/klines"
)

type Indicator interface {
	Mon() Monitor
	Period() int
	Active(klns1 []klines.Kline, klns2 []klines.Kline) (bool, error)
}

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
	mon        Monitor
	valPos     ValuePos
	bbLine     BBLine
	period     int
	multiplier float64
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
