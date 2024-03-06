package agent2

import (
	"fmt"
	"math"

	"github.com/varga-lp/data/klines"
)

type Trade struct {
	OpenTime     int64         `json:"open_time"`
	CloseTime    int64         `json:"close_time"`
	Reason       ClosingReason `json:"reason"`
	NetProfit    float64       `json:"net_profit"`
	DurationSecs int64         `json:"duration_scs"`
}

var (
	ErrPositionCantBeNilForTrade      = fmt.Errorf("position can't be nil when creating a trade")
	ErrPositionCloseIsNotLTECloseTime = fmt.Errorf("pos.Long.CloseTime should be less than equal to longClose.CloseTime")
	ErrBucketEndTimeIsNotGTStartTime  = fmt.Errorf("bucket end time is not gt than start time")
)

func NewTrade(pos *Position, cr ClosingReason, longClose klines.Kline, shortClose klines.Kline) (*Trade, error) {
	if pos == nil {
		return nil, ErrPositionCantBeNilForTrade
	}
	if pos.Long.CloseTime > longClose.CloseTime {
		return nil, ErrPositionCloseIsNotLTECloseTime
	}

	return &Trade{
		OpenTime:     pos.Long.OpenTime,
		CloseTime:    longClose.CloseTime,
		Reason:       cr,
		NetProfit:    pos.NetProfit(longClose, shortClose),
		DurationSecs: (longClose.CloseTime - pos.Long.CloseTime) / 1_000,
	}, nil
}

func (tr *Trade) String() string {
	return fmt.Sprintf("[trd] OT=%d, CT=%d, DSecs=%d, Net=%.2f, R=%s",
		tr.OpenTime, tr.CloseTime, tr.DurationSecs, tr.NetProfit, tr.Reason)
}

type Bucket struct {
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"`
	Trades    []*Trade `json:"trades"`
}

func NewBucket(startTime int64, endTime int64) (*Bucket, error) {
	if !(endTime > startTime) {
		return nil, ErrBucketEndTimeIsNotGTStartTime
	}

	return &Bucket{
		StartTime: startTime,
		EndTime:   endTime,
		Trades:    make([]*Trade, 0),
	}, nil
}

func (bu *Bucket) LastTrade() *Trade {
	tLen := len(bu.Trades)

	if tLen == 0 {
		return nil
	}
	return bu.Trades[tLen-1]
}

func (bu *Bucket) AppendTrade(pos *Position, cr ClosingReason, longClose klines.Kline, shortClose klines.Kline) error {
	trade, err := NewTrade(pos, cr, longClose, shortClose)
	if err != nil {
		return err
	}

	bu.Trades = append(bu.Trades, trade)
	return nil
}

const (
	dayLenMillis = float64(24 * 60 * 60 * 1_000)
)

func (bu *Bucket) HitRatio() float64 {
	if len(bu.Trades) == 0 {
		return 0.0
	}

	profitableTrades := 0.0
	for _, tr := range bu.Trades {
		if tr.NetProfit > 0.0 {
			profitableTrades++
		}
	}

	hr := profitableTrades / float64(len(bu.Trades))
	return roundTo4d(hr)
}

func (bu *Bucket) ProfitPerDay() float64 {
	days, totalProfit := float64(bu.EndTime-bu.StartTime)/dayLenMillis, 0.0
	for _, tr := range bu.Trades {
		totalProfit += tr.NetProfit
	}

	ppd := totalProfit / days
	return roundTo4d(ppd)
}

func roundTo4d(val float64) float64 {
	return math.Round(val*10_000.0) / 10_000.0
}
